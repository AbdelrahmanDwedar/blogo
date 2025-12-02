package main

import (
	"fmt"
	"log"
	"os"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
	"AbdelrahmanDwedar/blogo/internal/infrastructure/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	blogRepo := database.NewBlogRepository(db)

	fmt.Println("🌱 Seeding database...")

	// Create sample users
	users := []struct {
		username    string
		email       string
		displayName string
	}{
		{"alice", "alice@example.com", "Alice Wonder"},
		{"bob", "bob@example.com", "Bob Builder"},
		{"charlie", "charlie@example.com", "Charlie Brown"},
	}

	userIDs := []int64{}

	for _, u := range users {
		user := entity.NewUser(u.username, u.email, u.displayName)
		err := userRepo.Create(user)
		if err != nil {
			log.Printf("Warning: Could not create user %s: %v\n", u.username, err)
			continue
		}
		userIDs = append(userIDs, user.ID)
		fmt.Printf("✅ Created user: %s (ID: %d)\n", u.username, user.ID)
	}

	if len(userIDs) < 2 {
		log.Fatal("Not enough users created for seeding")
	}

	// Create sample blogs
	blogs := []struct {
		title       string
		description string
		body        string
		authorIdx   int
	}{
		{
			"Getting Started with Go",
			"A beginner's guide to Go programming",
			"Go is an open source programming language that makes it easy to build simple, reliable, and efficient software. In this comprehensive guide, we'll explore the fundamentals of Go programming...",
			0,
		},
		{
			"Building RESTful APIs",
			"How to build clean REST APIs with Go",
			"RESTful APIs are the backbone of modern web applications. In this post, I'll share best practices and patterns for building robust APIs using Go...",
			0,
		},
		{
			"Docker for Beginners",
			"Understanding containerization",
			"Docker has revolutionized how we deploy applications. Let's dive into what containers are and why they're so powerful...",
			1,
		},
		{
			"Database Design Tips",
			"Essential database design principles",
			"Good database design is crucial for application performance and maintainability. Here are some key principles to follow...",
			1,
		},
		{
			"The Art of Code Review",
			"How to give and receive effective code reviews",
			"Code reviews are an essential part of software development. They help catch bugs, share knowledge, and maintain code quality...",
			2,
		},
	}

	blogIDs := []int64{}

	for _, b := range blogs {
		if b.authorIdx >= len(userIDs) {
			continue
		}

		blog := entity.NewBlog(b.title, b.description, b.body, userIDs[b.authorIdx])
		err := blogRepo.Create(blog)
		if err != nil {
			log.Printf("Warning: Could not create blog '%s': %v\n", b.title, err)
			continue
		}
		blogIDs = append(blogIDs, blog.ID)
		fmt.Printf("✅ Created blog: %s (ID: %d)\n", b.title, blog.ID)
	}

	// Create some follow relationships
	if len(userIDs) >= 3 {
		// Alice follows Bob and Charlie
		userRepo.Follow(userIDs[0], userIDs[1])
		userRepo.Follow(userIDs[0], userIDs[2])

		// Bob follows Alice
		userRepo.Follow(userIDs[1], userIDs[0])

		// Charlie follows everyone
		userRepo.Follow(userIDs[2], userIDs[0])
		userRepo.Follow(userIDs[2], userIDs[1])

		fmt.Println("✅ Created follow relationships")
	}

	// Create some likes
	if len(blogIDs) >= 3 && len(userIDs) >= 3 {
		// Alice likes blogs 2 and 3
		blogRepo.Like(blogIDs[1], userIDs[0])
		blogRepo.Like(blogIDs[2], userIDs[0])

		// Bob likes blog 1 and 5
		blogRepo.Like(blogIDs[0], userIDs[1])
		if len(blogIDs) >= 5 {
			blogRepo.Like(blogIDs[4], userIDs[1])
		}

		// Charlie likes all blogs
		for _, blogID := range blogIDs {
			blogRepo.Like(blogID, userIDs[2])
		}

		fmt.Println("✅ Created likes")
	}

	fmt.Println("\n🎉 Database seeding complete!")
	fmt.Printf("\nCreated:\n")
	fmt.Printf("  - %d users\n", len(userIDs))
	fmt.Printf("  - %d blog posts\n", len(blogIDs))
	fmt.Println("\nSample credentials:")
	for i, u := range users {
		if i >= len(userIDs) {
			break
		}
		fmt.Printf("  - Username: %s, Email: %s\n", u.username, u.email)
	}

	os.Exit(0)
}


