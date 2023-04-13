const urlInputBox = document.querySelector('#input');
const generateUrlButton = document.querySelector('#generate-button');
const resultContainer = document.querySelector('.app__shorturl-result');
const actualUrl = document.querySelector('.app__actual-url');
const shortUrl = document.querySelector('.app__short-url');

generateUrlButton.onclick =  (event) => {
    const url = urlInputBox.value.trim();
    getShortUrl(url);
};

function getShortUrl(url) {
    fetch(`/getShortUrl`, {
        method: 'POST',
        headers: {
            'Content-Type': undefined
        },
        body: JSON.stringify({
            url
        })
    })
    .then((resp) => resp.json())
    .then((result) => {
        if (result.hasOwnProperty('Response')) {
            renderURLs(result['Response']);
        }
    })
    .catch((error) => {
        console.log(error);
    });
}

function renderURLs(response) {
    const {
        ActualURL, ShortURL
    } = response;
    resultContainer.style.display = 'initial';
    actualUrl.innerHTML = ActualURL;
    shortUrl.innerHTML = ShortURL;
}