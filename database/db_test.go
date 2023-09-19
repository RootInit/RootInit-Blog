package database

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestGetArticle(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}
	article, err := db.GetArticle(-1, false, true)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(article)
}

func TestGetComments(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}

	comments, err := db.GetArticleComments(1)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(comments)

	comment, err := db.GetComment(1)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(comment)
}

func TestGetUsers(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}

	users, err := db.GetUsers()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(users)

	user, err := db.GetUser(1)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(user)
}

func TestGetTags(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}

	tags, err := db.GetTags()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(tags)
}

func TestIttr(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}
	aIttr := db.GetArticleIttr()
	for aIttr.Next() {
		log.Println("test", aIttr.Article)
	}
	if aIttr.Error != nil {
		log.Println(aIttr.Error)
	}
}

func TestDataUsers(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		tag, err := db.NewUser(
			fmt.Sprintf("Username%d", i),
			fmt.Sprintf("Email%d@email.com", i),
			"",
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(tag)
	}
}

func TestDataTags(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		tag, err := db.NewTag(
			fmt.Sprintf("Tag%d", i),
			"asterisk",
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(tag)
	}
}

func TestDataCategories(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		tag, err := db.NewCategory(
			fmt.Sprintf("Category%d", i),
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(tag)
	}
}

func TestDataArticles(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}

	testArticleBody :=
		`<p>Once a month or so, I have the privilege of sitting down with Editor-in-Chief Elliot Williams
	to record the Hackaday Podcast. It’s a lot of fun spending a couple of hours geeking out
	together, and we invariably go off on ridiculous tangents with no chance of making the final
	cut, except perhaps as fodder for the intro and outro. It’s a lot of work, especially for
	Elliot, who has to edit the raw recordings, but it’s also a lot of fun.</p>
<p>Of course, we do the whole thing virtually, and we have a little ritual that we do at the
	start: the clapping. We take turns clapping our hands into our microphones three times, with
	the person on the other end of the line doing a clap of his own synchronized with the final
	clap. That gives Elliot an idea of how much lag there is on the line, which allows him to
	synchronize the two recordings. With him being in Germany and me in Idaho, the lag is pretty
	noticeable, at least a second or two.</p>
<p>Every time we perform this ritual, I can’t help but wonder about all the gear that makes it
	possible, including the fiber optic cables running beneath the Atlantic Ocean. Undersea
	communications cable stitch the world together, carrying more than 99% of transcontinental
	internet traffic. They’re full of fascinating engineering, but for my money, the inline
	optical repeaters that boost the signals along the way are the most interesting bits, even
	though — or perhaps especially because — they’re hidden away at the bottom of the sea.</p>
<h2>Better Than Coax</h2>
<p>Most of <a
		href="https://hackaday.com/2016/03/18/what-lies-beneath-the-first-transatlantic-communications-cables/">the
		long history of transoceanic communications</a> has been dominated by one material:
	copper. From the earliest telegraph cables right through to the coaxial cables carrying
	thousands of multiplexed telephone and television signals, copper conductors did the bulk of
	the work for almost all of the 20th century. That began to change in 1988 with the laying of
	the first transatlantic fiber-optic telephone cable, TAT-8. With a capacity of 40,000
	simultaneous phone calls on just two pairs of single-mode glass fibers (with one pair in
	reserve), TAT-8 bested the most advanced coaxial transatlantic cables by a factor of ten.
</p>
<pre><div class="codeBlock w75 center"><div><span>Code_Title.js</span><button class="hide">Copy</button></div>
			<code class="language-js">// a function to do a thing
document.addEventListener("DOMContentLoaded", function (event) {
(function () {
const navbar = document.querySelector(".navbar");

// Navbar Scroll Hiding
let prevScrollPos = window.scrollY;
function showHideNavbar() {
	const currentScrollPos = window.scrollY;
	if (prevScrollPos < currentScrollPos) {
		navbar.classList.add("navbar-show");
		navbar.classList.remove("navbar-hide");
	} else {
		navbar.classList.remove("navbar-show");
		navbar.classList.add("navbar-hide");
	}
	prevScrollPos = currentScrollPos;
}
window.addEventListener("scroll", showHideNavbar);
})();
	</code></div></pre>
`
	for i := 0; i < 50; i++ {
		title := fmt.Sprintf("Test Article %d", i)
		article, err := db.NewArticle(
			title,
			getUrlTitle(title),
			fmt.Sprintf("Test Description %d", i),
			1,
			time.Now(),
			testArticleBody,
			"testThumb.png",
			[]int{1, 2, 3},
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(article)
	}
}

func TestDataComments(t *testing.T) {
	db, err := InitializeDB("../blog.db")
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < 50; i++ {
		comment, err := db.NewComment(
			i,
			0,
			1,
			time.Now(),
			fmt.Sprintf("Test Article %d", i),
		)
		if err != nil {
			log.Println(err)
		}
		log.Println(comment)
		for ii := 1; ii < 4; ii++ {
			comment, err := db.NewComment(
				i,
				comment.Id,
				ii,
				time.Now(),
				fmt.Sprintf("Test Article %d", i),
			)
			if err != nil {
				log.Println(err)
			}
			log.Println(comment)
		}
		if err != nil {
			log.Println(err)
		}
		log.Println(comment)

	}
}
