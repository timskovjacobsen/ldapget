package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/timskovjacobsen/ldapget/client"
	"github.com/timskovjacobsen/ldapget/config"
	"github.com/timskovjacobsen/ldapget/layout"
)

var cfg *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ldapget",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		globalConfig, err := config.ReadConfig("")
		if err != nil {
			return fmt.Errorf("failed to read config: %v", err)
		}
		cfg = globalConfig // config becomes globally available
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bare application")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func GroupsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groups",
		Short: "List AD groups with related information",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			groups := client.Groups(cfg)

			var headerStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#cbba82")).
				BorderTop(true).
				BorderBottom(true)
			var nameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cbba82"))

			fmt.Println(headerStyle.Render("\nAD Groups Information:"))
			separator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#555555"))

			width, _, _ := term.GetSize(int(os.Stderr.Fd()))
			horizontalRule := separator.Render(strings.Repeat("─", width))
			fmt.Println(horizontalRule)

			for i, group := range groups {
				fmt.Printf("%d. ", i+1)
				fmt.Println(nameStyle.Render(fmt.Sprintf("%s", group.Name)))
				fmt.Printf("   🗺️ %s\n", group.DN)
				fmt.Printf("   🏷️ %s group\n", group.Kind)
				if group.SystemCreated {
					fmt.Printf("   System created: %s\n", "yes")
				}
				fmt.Printf("   🎯 %s scope\n", group.Scope)
				if group.Description != "" {
					fmt.Printf("   📝 %s\n", group.Description)
				}
				fmt.Printf("   👥 %d members\n", group.Members)
				fmt.Println(horizontalRule)
			}
		},
	}
	return cmd
}

func GroupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "group short",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			group := args[0]
			baseDN := cfg.Client.Search.RootDN
			conn, err := client.BindToLdapServer(*cfg)
			if err != nil {
				log.Fatalf("failed to bind to ldap server: %v", err)
			}
			result, err := client.Group(conn, baseDN, group)

			if err != nil {
				log.Fatalf("Failed to search LDAP server for group: %v", err)
			}
			if len(result.Entries) == 0 {
				log.Fatalf("Group not found")
			}

			var groupList []string
			for _, entry := range result.Entries {
				groupList = append(groupList, entry.GetAttributeValue("cn"))
			}
			slices.Sort(groupList)

			enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#646464"))
			itemStyle := lipgloss.NewStyle().MarginLeft(1)
			formattedList := list.New(groupList).ItemStyle(itemStyle).EnumeratorStyle(enumeratorStyle).Enumerator(layout.Arabic)
			fmt.Println(formattedList)

		},
	}
	return cmd
}

func UserCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user <USER_INITIALS>",
		Short: "Look up the groups that a user is member of",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			user := args[0]
			baseDN := cfg.Client.Search.RootDN
			conn, err := client.BindToLdapServer(*cfg)
			if err != nil {
				log.Fatalf("failed to bind to ldap server: %v", err)
			}
			defer conn.Close()
			result, err := client.User(conn, baseDN, user)
			if err != nil {
				log.Fatalf("Failed to search LDAP server for user: %v", err)
			}
			if len(result.Entries) == 0 {
				log.Fatalf("User not found")
			}

			// Print out the groups
			for _, entry := range result.Entries {
				DNFields := strings.Split(entry.DN, ",")
				var name string
				for _, field := range DNFields {
					if strings.HasPrefix(field, "CN=") {
						name = field[3:]
					}
				}
				var headerStyle = lipgloss.NewStyle().
					BorderStyle(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#cbba82")).
					BorderTop(true).
					BorderBottom(true)
				mail := entry.GetAttributeValue("mail")
				fmt.Println(headerStyle.Render(fmt.Sprintf("%s (%s)", name, mail)))
				itemStyle := lipgloss.NewStyle().MarginLeft(1)
				enumeratorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#646464"))

				groupList := list.New().ItemStyle(itemStyle).EnumeratorStyle(enumeratorStyle).Enumerator(layout.Arabic)
				for _, attr := range entry.GetAttributeValues("memberOf") {
					fields := strings.Split(attr, ",")
					for _, field := range fields {
						if strings.HasPrefix(field, "CN=") {
							group := field[3:]
							groupList.Item(group)
						}
					}
				}
				fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#cbba82")).Render("AD GROUPS"))
				fmt.Println(groupList)
			}
		},
	}
	return cmd
}

func init() {
	rootCmd.AddCommand(GroupsCommand())
	rootCmd.AddCommand(GroupCommand())
	rootCmd.AddCommand(UserCommand())
}
