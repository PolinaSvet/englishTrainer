
function homePage() {
    window.location.href = "/";
}

function cardsPage() {
    window.location.href = "/cards";
}

function loginPage() {
    window.location.href = "/login";
}

function registrationPage() {
    window.location.href = "/registration";
}

// Функция для воспроизведения текста
function speakText(text) {
    if ('speechSynthesis' in window) {
        const utterance = new SpeechSynthesisUtterance(text);

        utterance.lang = 'en-US'; // Устанавливаем язык (английский)
        utterance.rate = 1; // Скорость речи (1 - нормальная)
        utterance.pitch = 1; // Высота голоса (1 - нормальная)
        utterance.volume = 1; // Громкость (1 - максимальная)

        window.speechSynthesis.speak(utterance);
    } else {
        showAlarm('Ваш браузер не поддерживает синтез речи.');
    }
}

// Функция для вывода alarm mess
function showAlarm(text) {
    
    if (String(text) == "undefined") {
        return
    }

    console.log(text);

    const alarm = document.createElement("div");
    alarm.className = "alarm";

    // Добавляем текст сообщения
    const message = document.createElement("span");
    message.textContent = JSON.stringify(text);
    alarm.appendChild(message);

    // Создаем кнопку закрытия
    const closeButton = document.createElement("button");
    closeButton.textContent = "×"; // Символ крестика
    closeButton.onclick = function () {
        // Удаляем сообщение при нажатии на кнопку
        alarm.remove();
    };
    alarm.appendChild(closeButton);

    // Добавляем сообщение в контейнер
    const container = document.getElementById("alarm-container");
    container.appendChild(alarm);

    // Автоматическое удаление через 5 секунд (опционально)
    //setTimeout(() => {
    //    if (alarm.parentElement) {
    //        alarm.remove();
    //    }
    //}, 5000);
}
