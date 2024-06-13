function toggleForms() {
    var loginForm = document.getElementById('loginForm');
    var registerForm = document.getElementById('registerForm');
    if (loginForm.style.display === "none") {
        loginForm.classList.remove('fade-out');
        loginForm.classList.add('fade-in');
        registerForm.classList.remove('fade-in');
        registerForm.classList.add('fade-out');
        setTimeout(function () {
            loginForm.style.display = "block";
            registerForm.style.display = "none";
        }, 500);
    } else {
        loginForm.classList.remove('fade-in');
        loginForm.classList.add('fade-out');
        registerForm.classList.remove('fade-out');
        registerForm.classList.add('fade-in');
        setTimeout(function () {
            loginForm.style.display = "none";
            registerForm.style.display = "block";
        }, 500);
    }
}


var registerFormID = document.getElementById("registerFormID")
registerFormID.addEventListener("submit", async (event) => {
    event.preventDefault()
    const data = getDataForm(registerFormID)
    data.age = strToInt(data.age)

    try {
        const urlRegister = 'http://localhost:8080/register'
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        };
        const response = await fetch(urlRegister, requestOptions)

        if (!response.ok) {
            throw new Error(`HTTP error status: ${response.status}`);
        }

        const result = await response.json();
        console.log(result);

    } catch (error) {
        console.error('Erreur lors de l\'envoi des données:', error);
    }


})

function getDataForm(form) {
    const dataForm = new FormData(form)
    var data = Object.fromEntries(dataForm.entries())
    return data
}

function strToInt(str) {
    num = Number(str)
    return parseInt(num, 10)
}

//********************************************************************************************************************** */

var loginFormID = document.getElementById("loginFormID")
loginFormID.addEventListener("submit", async (event) => {
    event.preventDefault()

    const data = getDataForm(loginFormID)

    try {
        const urlLogin = 'http://localhost:8080/login'
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        };
        const response = await fetch(urlLogin, requestOptions)

        if (!response.ok) {
            throw new Error(`HTTP error status: ${response.status}`);
        }

        const result = await response.json();
        console.log(result);

    } catch (error) {
        console.error('Erreur lors de l\'envoi des données:', error);
    }

})