package main

import (
    "context"
    "log"
    "strings"
    "time"
	"os/exec"
	"fmt"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/chromedp/chromedp"
)

type Company struct {
	Name	string
	Link	string
}

type MoreCompaniesMsg struct {
	Category string
	Companies []Company
}

type LoadingMsg struct {
	Category string
	Loading bool
}

type model struct {
    companies			[]Company
    categories			[]string
    selectedCompany  	int
    width     			int
    height    			int
	navbar				bool
	selectedCategory 	int
	cache 				*CacheData
	isLoading			bool
	loadingFrame		int
}


func scrapeAllData() (*CacheData, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var homeLinks []string
	err := chromedp.Run(ctx, 
		chromedp.Navigate("https://startups.gallery"),
		chromedp.Sleep(5),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &homeLinks),
	)
	if err != nil {
		return nil, err
	}
	homeCompanies := extractCompanies(homeLinks)
	categoryLinks := getCategoryLinks(homeLinks)

	categoryLinks["All"] = "https://startups.gallery"  

	var categories []string
	categories = append(categories, "All")
	categoryData := make(map[string][]Company)
	
	categoryData["All"] = homeCompanies
	for name, link := range categoryLinks {
	
		var nameLinks []string
		err := chromedp.Run(ctx, 
			chromedp.Navigate(link),
			chromedp.Sleep(5),
			chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &nameLinks),
		)

		if err != nil {
			return nil, err
		}

		var nameCompanies = extractCompanies(nameLinks)
		categoryData[name] = nameCompanies
	}
	return &CacheData {
		Date: time.Now().Format("2006-01-01"),
		CategoryCompanyMap: categoryData,
		CategoriesLinkMap: categoryLinks,
	}, nil
}

func getCategoryLinks(allLinks []string) map[string]string {
	categoryMap := make(map[string]string)

	for _, link := range allLinks {
		if strings.Contains(link, "/categories/") {
            parts := strings.Split(link, "/")
			categoryMap[parts[5]] = link
		}
	}
	return categoryMap
}

func openInFirefox(url string) {
	if err := exec.Command("firefox", url).Start(); err == nil {
		return
	}
}

func getData() (*CacheData, error) {
	cacheFilePath := GetCacheFilePath()

	if IsCacheValid(cacheFilePath) {
        fmt.Println("📁 Loading from cache...")
		cache, err := LoadCache(cacheFilePath)
		if err == nil {
			fmt.Printf("✓ Loaded %d categories from cache\n", len(cache.CategoriesLinkMap))
			return cache, nil
		}
        fmt.Println("⚠️ Cache read failed, fetching fresh data...")
	} else {
        fmt.Println("🌐 Cache expired, fetching fresh data...")
	}
	cache, err := scrapeAllData()
	if err != nil {
		return nil, err
	}

    if err := SaveData(cacheFilePath, cache); err != nil {
        fmt.Printf("⚠️ Failed to save cache: %v\n", err)
    } else {
        fmt.Println("💾 Data cached successfully!")
    }
	return cache, nil
}

func (m *model) SwitchCategory(categoryName string) {
	if companies, exists := m.cache.CategoryCompanyMap[categoryName]; exists {
		m.companies = companies
		m.selectedCompany = 0
	}
}

func initialModel(cache *CacheData) model {
	companies := cache.CategoryCompanyMap["All"]	
    desiredOrder := []string{"All", "ai", "remote", "fintech", "devtools", "design", 
                            "productivity", "analytics", "healthcare", "seed", "pre-seed", 
                            "series-a", "series-b", "series-c"}
	var categories []string
	for _, categoryName := range desiredOrder {
		if _, exists := cache.CategoryCompanyMap[categoryName]; exists {
			categories = append(categories, categoryName)
		}
	}
	return model {
        companies: companies,
        categories:	categories,
        selectedCompany: 0,
		selectedCategory: 0,
		cache: cache,
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
			if m.navbar {
				m.selectedCategory = (m.selectedCategory - 1 + len(m.categories)) % len(m.categories)  
			} else {
				m.selectedCompany = (m.selectedCompany - 1 + len(m.companies)) % len(m.companies)
			}
        case "right", "l":
			if m.navbar {
				m.selectedCategory = (m.selectedCategory + 1 + len(m.categories)) % len(m.categories)  
			} else {
				m.selectedCompany = (m.selectedCompany + 1 + len(m.companies)) % len(m.companies)
			}
        case "up", "k":
			if !m.navbar {
				m.selectedCompany = (m.selectedCompany - 3 + len(m.companies)) % len(m.companies)
			} 			
        case "down", "j":
			if !m.navbar {
				if m.selectedCompany >= len(m.companies) - 3 && !m.isLoading && m.categories[m.selectedCategory] != "All"{
					return m, func() tea.Msg {
						return LoadingMsg {Category: m.categories[m.selectedCategory], Loading: true}
					}
				}
				m.selectedCompany = (m.selectedCompany + 3 + len(m.companies)) % len(m.companies) 				
			}

        case "enter", "o":
			if m.navbar {
				m.SwitchCategory(m.categories[m.selectedCategory])
				m.navbar = false
			} else {
				url := m.companies[m.selectedCompany].Link
				go openInFirefox(url)
			}

		case "n":
			m.navbar = !m.navbar

        }

	case LoadingMsg:
		if msg.Category == m.categories[m.selectedCategory] {
			m.isLoading = msg.Loading
			if msg.Loading {
				return m, func() tea.Msg {
					companies := ScrapeMoreCompanies(msg.Category, len(m.companies))
					return MoreCompaniesMsg {Category: msg.Category, Companies: companies} 
				}
			}
		}

	case MoreCompaniesMsg:
		if msg.Category == m.categories[m.selectedCategory] {
			m.companies = append(m.companies, msg.Companies...)
			m.cache.CategoryCompanyMap[msg.Category] = m.companies
			m.isLoading = false
		} 
    }
    return m, nil
}

func ScrapeMoreCompanies(category string, currentCount int) []Company {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()
    
    var allLinks []string
    err := chromedp.Run(ctx,
        chromedp.Navigate(fmt.Sprintf("https://startups.gallery/categories/%s", category)),
        chromedp.Sleep(2*time.Second),
        
        // Simulate scrolling to load more
        chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
        chromedp.Sleep(3*time.Second),
        
        chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &allLinks),
    )
    
    if err != nil {
        return []Company{} // Return empty on error
    }
    
    allCompanies := extractCompanies(allLinks)
    
    // Return only the new companies (skip the ones we already have)
    if len(allCompanies) > currentCount {
        return allCompanies[currentCount:]
    }
    return []Company{}
}

func (m model) View() string {
	var navbarStyle, selectedNavbarStyle lipgloss.Style

	if m.navbar {
		navbarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("12")).
			Padding(0, 1).
			MarginBottom(1)
		
		selectedNavbarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("5")).
			Padding(0, 1).
			MarginBottom(1)
	} else {
		navbarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Background(lipgloss.Color("0")).
			Padding(0, 1).
			MarginBottom(1)

		selectedNavbarStyle = navbarStyle.Copy().
			Foreground(lipgloss.Color("15"))
	}

	maxNavBarWidth := m.width + 10
	var visibleNavItems []string
	currentWidth := 0
	startIdx := max(0, m.selectedCategory - 3)
	endIdx := min(len(m.categories), startIdx + 6)


	leftArrow := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render("◀")
    visibleNavItems = append(visibleNavItems, leftArrow)
    currentWidth += 2
	for i := startIdx; i < endIdx; i++ {
		category := m.categories[i]
		itemWidth := len(category) + 4

		if currentWidth + itemWidth <= maxNavBarWidth {

			var style lipgloss.Style	
			if i == m.selectedCategory {
				style = selectedNavbarStyle
			} else {
				style = navbarStyle
			}
			visibleNavItems = append(visibleNavItems, style.Render(category))
			currentWidth += itemWidth
		} else {
			break
		}		
	}
	rightArrow := lipgloss.NewStyle().
        Foreground(lipgloss.Color("8")).
        Render("▶")
    visibleNavItems = append(visibleNavItems, rightArrow)

	navbar := lipgloss.JoinHorizontal(lipgloss.Top, visibleNavItems...)

    cardStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("8")).
        Padding(1).
        Width(25).
        Height(3)

    selectedCardStyle := cardStyle.Copy().
        BorderForeground(lipgloss.Color("10")).
        Background(lipgloss.Color("0")).
		Bold(true)

    var rows []string
    for i := 0; i < len(m.companies); i += 3 {
        var row []string

        for j := 0; j < 3 && i+j < len(m.companies); j++ {
            idx := i + j
            company := m.companies[idx]
            var card string
            if idx == m.selectedCompany {
                card = selectedCardStyle.Render(fmt.Sprintf("▶ %s", company.Name))
            } else {
                card = cardStyle.Render(company.Name)
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
	var instructions string
	
	if m.navbar {
        instructions = "NAVBAR MODE: ←→ to select category • Enter to switch • n to exit navbar"
	} else {
		loadingBar := ""
		if m.isLoading {
			loadingBar = " • " + renderProgressBar(20)
		}
        instructions = fmt.Sprintf(
			"Category: %s (%d companies)%s • n for navbar • arrows to navigate • o/enter to open • q to quit",
			m.categories[m.selectedCategory], len(m.companies), loadingBar)
	}

    instructionStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("8")).
        MarginTop(2)

    content := lipgloss.JoinVertical(lipgloss.Left,navbar,grid,instructionStyle.Render(instructions))

    return content
}

func renderProgressBar(width int) string {
    bar := strings.Repeat("█", width/4) + strings.Repeat("░", width/4)
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("12")).
        Render(bar)
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
					Name: companyName,
					Link: link,
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
	cache, err := getData()
	if err != nil {
		log.Fatal(err)
	}
	p := tea.NewProgram(initialModel(cache), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
