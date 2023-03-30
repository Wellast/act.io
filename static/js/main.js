const redirect = (path) => window.location.href = path;
const getToken = () => localStorage.getItem('token');
const boundClickToFunc = (elementId, func) => document.getElementById(elementId).addEventListener("click", func);
const isAuthorized = () => getToken() ? true : false;

const changeBackground = () => {
    const base64ToBlob = (url) => fetch(`${url}`).then(res => res.blob());
    const blobToBase64 = (blob) => new Promise((res) => {
        const reader = new FileReader();
        reader.onloadend = () => res(reader.result);
        reader.readAsDataURL(blob);
    });

    //  setup cached image
    const bs64Image = localStorage.getItem('image');
    base64ToBlob(bs64Image)
        .then((imageBlob) => {
            const imageObjectURL = URL.createObjectURL(imageBlob);
            document.getElementById('background').src = imageObjectURL;
        });

    //  fetch a new image
    const cancelFunc = progressCreate();
    const tags = (localStorage.getItem('tags') || 'wallpaper');
    const url = `https://source.unsplash.com/random/${window.screen.width}x${window.screen.height}?${tags}`;
    const controller = new AbortController();
    setTimeout(() => controller.abort(), 5000);
    fetch(url, { signal: controller.signal })
    .then((res) => res.blob())
    .then((imageBlob) => {
        const imageObjectURL = URL.createObjectURL(imageBlob);
        document.getElementById('background').src = imageObjectURL;
        blobToBase64(imageBlob)
        .then((result) => {
            localStorage.setItem('image', result);
            cancelFunc();
        });
    })
    .catch(cancelFunc)
}

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
            .catch(console.warn)
    } else {
        showLoginButton();
    }

    changeBackground();
})();
