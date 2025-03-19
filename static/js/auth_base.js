
// Функция для обновления меню
function updateAuthMenu() {
    const objData = localStorage.getItem("user");
    const obj = objData ? JSON.parse(objData) : null;

    const authLink = document.getElementById("auth-link");
    const authUser = document.getElementById("auth-user");
    const authMail = document.getElementById("auth-mail");
    const authLogin = document.getElementById("auth-login");
    const authRegistration = document.getElementById("auth-registration");
    const authLogout = document.getElementById("auth-logout");

    if (obj) {
        // Если пользователь авторизован
        authLink.textContent = obj.user.name ? obj.user.name : "Login";
        authUser.textContent = obj.user.name ? obj.user.name : "Name";
        authMail.textContent = obj.user.mail ? obj.user.mail : "Mail";

        // Показываем Logout, скрываем Login и Registration
        authLogin.style.display = "none";
        authRegistration.style.display = "none";
        authLogout.style.display = "block";
    } else {
        // Если пользователь не авторизован
        authLink.textContent = "Login";
        authUser.textContent = "Name";
        authMail.textContent = "Mail";

        // Показываем Login и Registration, скрываем Logout
        authLogin.style.display = "block";
        authRegistration.style.display = "block";
        authLogout.style.display = "none";
    }
}

// Инициализация меню при загрузке страницы
document.addEventListener("DOMContentLoaded", updateAuthMenu);

