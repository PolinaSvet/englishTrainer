document.getElementById("authForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const formData = {
        name: document.getElementById("name")?.value,
        mail: document.getElementById("mail")?.value,
        telega: document.getElementById("telega")?.value,
        password: document.getElementById("password")?.value,
    };


    let endpoint = "/auth/logout";
    let method = "DELETE";
    let authTitel = document.getElementById("authTitel")?.textContent;

    console.log(":"+authTitel+":")

    if (authTitel === "Login") {
        endpoint = "/auth/login";
        method = "PUT";
    } else if (authTitel === "Register") {
        endpoint = "/auth/register";
        method = "POST";
    }
    console.log(endpoint,method)

    const response = await fetch(endpoint, {
        method: method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(formData),
    });

    

    const authError = document.getElementById("authError");
    if (response.ok) {
        const user = await response.json();
        console.log(user)
        localStorage.setItem("user", JSON.stringify(user)); // Сохраняем данные пользователя
        authError.textContent = "Success!";
        authError.classList.remove("failure");
        authError.classList.add("success");
        if (method === "DELETE") localStorage.removeItem("user");
        if (method === "POST") localStorage.removeItem("user");
        updateAuthMenu(); 
        homePage();
    } else {
        const error = await response.text();
        authError.textContent = error;
        authError.classList.remove("success");
        authError.classList.add("failure");
    }
});

