package main

import (
    "context"
    "log"
    "strings"
    "time"
	"os/exec"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/chromedp/chromedp"
)

type Company struct {
	name string
	link string
}

type model struct {
    companies  []Company
    categories []string
    selected   int
    width      int
    height     int
}

func openInFirefox(url string) {
	if err := exec.Command("firefox", url).Start(); err == nil {
		return
	}
}

func initialModel(companies []Company, categories []string) model {
    return model{
        companies:  companies,
        categories: categories,
        selected:   0,
    }
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        return m, nil

    case tea.KeyMsg:
        switch msg.String() {
        case "q":
            return m, tea.Quit
        case "left", "h":
            if m.selected > 0 {
                m.selected--
            }
        case "right", "l":
            if m.selected < len(m.companies)-1 {
                m.selected++
            }
        case "up", "k":
            if m.selected >= 3 {
                m.selected -= 3
            }
        case "down", "j":
            if m.selected+3 < len(m.companies) {
                m.selected += 3
            }
        case "enter", "o":
			url := m.companies[m.selected].link
    		go openInFirefox(url)
        }
    }
    return m, nil
}

func (m model) View() string {
    navbarStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("15")).
        Background(lipgloss.Color("8")).
        Padding(0, 1).
        MarginBottom(1)

    cardStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("8")).
        Padding(1).
        Width(25).
        Height(3)

    selectedCardStyle := cardStyle.Copy().
        BorderForeground(lipgloss.Color("12")).
        Background(lipgloss.Color("0"))

    var navItems []string
    for _, category := range m.categories {
        navItems = append(navItems, navbarStyle.Render(category))
    }
    navbar := lipgloss.JoinHorizontal(lipgloss.Top, navItems...)

    var rows []string
    for i := 0; i < len(m.companies); i += 3 {
        var row []string

        for j := 0; j < 3 && i+j < len(m.companies); j++ {
            idx := i + j
            company := m.companies[idx]
            var card string
            if idx == m.selected {
                card = selectedCardStyle.Render(company.name)
            } else {
                card = cardStyle.Render(company.name)
            }
            row = append(row, card)
        }
        
        // Fill empty columns if needed
        for len(row) < 3 {
            row = append(row, cardStyle.Render(""))
        }

        rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
    }
    grid := lipgloss.JoinVertical(lipgloss.Left, rows...)

    instructions := lipgloss.NewStyle().
        Foreground(lipgloss.Color("8")).
        MarginTop(2).
        Render("use arrow or vim keys for nav / q to quit")

    content := lipgloss.JoinVertical(
        lipgloss.Left,
        navbar,
        grid,
        instructions,
    )

    return content
}

func extractCompanies(links []string) []Company {
    var companies []Company
    for _, link := range links {
        if strings.Contains(link, "/companies/") {
            parts := strings.Split(link, "/")
            if len(parts) > 0 {
                companyName := parts[len(parts)-1]
                companyName = strings.ReplaceAll(companyName, "-", " ")
                companyName = strings.Title(companyName)
  				var c = Company {
					name: companyName,
					link: link,
				}
				companies = append(companies, c)
            }
        }
    }
    return companies
}

func extractCategories(links []string) []string {
    categoryMap := make(map[string]bool)

    for _, link := range links {
        if strings.Contains(link, "/categories/") {
            parts := strings.Split(link, "/")
            if len(parts) >= 5 && parts[3] == "categories" {
                category := parts[len(parts)-1]
                category = strings.ReplaceAll(category, "-", " ")
                category = strings.Title(category)
                categoryMap[category] = true
            }
        }
    }

    // Convert map to slice
    var categories []string
    categories = append(categories, "All") // Always add "All" first
    for category := range categoryMap {
        categories = append(categories, category)
    }

    return categories
}

func main() {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    var allLinks []string

    err := chromedp.Run(ctx,
        chromedp.Navigate("https://startups.gallery"),
        chromedp.Sleep(5*time.Second),
        chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &allLinks),
    )

    if err != nil {
        log.Fatal(err)
    }

    companies := extractCompanies(allLinks)
    categories := extractCategories(allLinks)

    p := tea.NewProgram(initialModel(companies, categories), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
