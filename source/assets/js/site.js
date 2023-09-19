//navbar show/hide
try {
(function () {
    const navbar = document.querySelector(".navBar");
    let prevScrollPos = window.scrollY;
    function showHideNavbar() {
        const currentScrollPos = window.scrollY;
        if (prevScrollPos < currentScrollPos) {
            navbar.classList.add("navbarShow");
            navbar.classList.remove("navbarHide");
        } else {
            navbar.classList.remove("navbarShow");
            navbar.classList.add("navbarHide");
        }
        prevScrollPos = currentScrollPos;
    }
    window.addEventListener("scroll", showHideNavbar);
})();
} catch (error) {
    console.error(error)
}

// Color Mode Switch
try{
(function () {
    const bulbButton = document.getElementById("colorModeToggle");
    const bulbButtonGlow = bulbButton.querySelector("#glowBg");

    // Check for stored cookie
    let darkModeEnabled = localStorage.getItem("darkModeEnabled");
    // Fallback to system color scheme
    if (darkModeEnabled == null) {
        if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
            darkModeEnabled = true;
        } else {
            darkModeEnabled = false;
        }
    }
    // Set initially
    if (darkModeEnabled == true) {
        setDarkMode();
    } else {
        setLightMode();
    }
    bulbButton.addEventListener("click", function () {
        if (darkModeEnabled == true) {
            setLightMode();
        } else {
            setDarkMode();
        }
    });
    function setLightMode() {
        document.documentElement.setAttribute("colorMode", "light");
        localStorage.setItem("darkModeEnabled", 0);
        bulbButtonGlow.style.opacity = 1;
        darkModeEnabled = false;
    }
    function setDarkMode() {
        document.documentElement.setAttribute("colorMode", "dark");
        localStorage.setItem("darkModeEnabled", 1);
        bulbButtonGlow.style.opacity = 0;
        darkModeEnabled = true;
    }
})();
} catch (error) {
    console.error(error)
}
