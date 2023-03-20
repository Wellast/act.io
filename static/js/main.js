const changeBackground = () => document.getElementById('background').src = `https://source.unsplash.com/random/${window.screen.width}x${window.screen.height}?wallpaper`;
const redirect = (path) => window.location.href = path;
const getToken = () => localStorage.getItem('token');
const boundClickToFunc = (elementId, func) => document.getElementById(elementId).addEventListener("click", func);
const isAuthorized = () => getToken() ? true : false;

const showLoginButton = () => {
    const retrieveAuthUrl = () => fetch('/auth')
        .then((res) => res.text())
        .then((redirectUrl) => {
            window.open(redirectUrl, "_self").focus();
        })
        .catch(err => alert(err))

    document.getElementById('login-button').hidden = false;
    boundClickToFunc('login-button', retrieveAuthUrl);
}

const showLogoutButton = () => {
    const logoutUser = () => {
        localStorage.removeItem('token');
        window.location.reload();
    }

    document.getElementById('logout-button').hidden = false;
    boundClickToFunc('logout-button', logoutUser);
}

const isLoginAttempt = () => {
    token = new URLSearchParams(window.location.search)
        .get('token');
    if (token) {
        localStorage.setItem('token', token);
        return true;
    }
    return false;
}


;(() => {
    isLoginAttempt() && redirect('/');

    if (isAuthorized()) {
        showLogoutButton();
        retrieveCalendarEvents(getToken())
            .then(drawEvents)
    } else {
        showLoginButton();
    }

    changeBackground();
})();
