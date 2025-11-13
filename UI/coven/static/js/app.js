function themeToggle_addThemeSwitchHandler() {
    const themeToggle = document.getElementById('themeToggle')
    const savedTheme = localStorage.getItem('bsTheme') || 'light'

    document.documentElement.setAttribute('data-bs-theme', savedTheme)

    if (!themeToggle){
        console.error('failed to init theme toggle')
        return;
    }

    themeToggle.checked = savedTheme === 'dark'

    themeToggle.addEventListener('change', function () {
        const newTheme = this.checked ? 'dark' : 'light'
        document.documentElement.setAttribute('data-bs-theme', newTheme)
        localStorage.setItem('bsTheme', newTheme)
    });
}

function mainMenu_addSelectedItemHandler() {
    const mainMenu = document.getElementById('main-menu')
    if (!mainMenu) {
        console.error('failed to get main-menu component')
        return
    }

    const menuIntems = mainMenu.querySelectorAll('#main-menu ul li')

    menuIntems.forEach(item => {
        const activeClassName = 'active'
        item.addEventListener('click', function(){
            menuIntems.forEach(li =>{
                li.classList.remove(activeClassName)
            })

            this.classList.add(activeClassName)
        })
    })
}