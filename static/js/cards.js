let currentCard = null;
let currentCardIndex = null;
let currentCardsData = null;
let currentScoreGuess = null;
let currentScoreAll = null;
let currentFinish = null;
let currentIsSart = null;

function initGame() {

    document.getElementById("card-container").style.display = "none";
    document.getElementById("score-container").style.display = "none";
    document.getElementById("words-table-container").style.display = "block";
    document.getElementById("lets-container").style.display = "block";
    document.getElementById("congrats-container").style.display = "none";

    document.getElementById("congrats-information").innerText = "";

    currentCardsData = JSON.parse(JSON.stringify(window.cardsData));
    currentScoreGuess = window.scoreGuess;
    currentScoreAll = window.scoreAll;
    currentFinish = window.finish;
    currentIsSart = false;

    if (currentCardsData == null){
        showAlarm(String(window.message))
        document.getElementById("lets-container").style.display = "none";
        document.getElementById("words-table-container").style.display = "none";
        return;
    }
   
    updateScore();
    updateWordsList();

}

function letsGame() {

    document.getElementById("card-container").style.display = "block";
    document.getElementById("score-container").style.display = "block";
    document.getElementById("words-table-container").style.display = "block";
    document.getElementById("lets-container").style.display = "none";
    document.getElementById("congrats-container").style.display = "none";

    currentIsSart = true;
    updateUI();
    
}

function finishGame() {
    
    document.getElementById("card-container").style.display = "none";
    document.getElementById("score-container").style.display = "block";
    document.getElementById("words-table-container").style.display = "block";
    document.getElementById("lets-container").style.display = "none";
    document.getElementById("congrats-container").style.display = "block";

    document.getElementById("congrats-information").innerText = "Congratulations! You have finished the game!";

    currentScoreGuess = currentScoreAll;
    currentIsSart = false
    currentFinish = true;
    updateScore();
    updateWordsList();
}


function updateWordsList() {
    
    let tbody = document.querySelector("#words-table tbody");
    tbody.innerHTML = "";

    document.getElementById("words-table-column3").style.display = currentIsSart ? 'none' : 'block';
    
    currentCardsData.forEach(card => {
        let row = document.createElement("tr");

        row.innerHTML = `
            <td>${card.word}</td>
            <td>${card.transcription}</td>
            <td style="display: ${currentIsSart ? 'none' : 'block'};">${card.translation}</td>
            <td onclick="speakText('${card.word}')" style="cursor: pointer;">&#9654;</td>
            <td>${card.attempt}</td>
            <td><div class="words-table-circle ${card.guess ? 'green' : 'grey'}"></div></td>
        `;

        tbody.appendChild(row);
    });
}

function updateScore() {
    
    document.getElementById("score-value").innerText = "Score: " + currentScoreGuess + "/" + currentScoreAll;
}

function getRandomUnGuessedCard() {
    let unguessedCards = currentCardsData.filter(card => !card.guess);

    if (unguessedCards.length === 0) return null;
    currentCardIndex = Math.floor(Math.random() * unguessedCards.length);
    return unguessedCards[currentCardIndex];
}

function updateUI() {

    currentCard = getRandomUnGuessedCard();
   
    if (!currentCard) {
        finishGame();
        return;
    }

    updateWordsList();
    updateScore();

    document.getElementById("card__word").innerText = currentCard.word;
    document.getElementById("card__word-back").innerText = currentCard.word;
    document.getElementById("card__transcription").innerText = currentCard.transcription;
    document.getElementById("card__translation").innerText = currentCard.translation;
    

    let buttons = document.querySelectorAll(".answer-button");
    buttons.forEach((btn, index) => {
        btn.innerText = currentCard.answers[index];
        btn.onclick = () => checkAnswer(currentCard, currentCard.answers[index]);
    });
   
    highlightActiveRow();
    
  
}

function highlightActiveRow() {
    let rows = document.querySelectorAll("#words-table tbody tr");
    rows.forEach(row => {
        row.classList.remove("active");
        if (row.cells[0].innerText === currentCard.word) {
            row.classList.add("active");
        }
    });
}

function sortTable(column) {
    currentCardsData.sort((a, b) => {
        if (column === "word") {
            return a.word.localeCompare(b.word);
        } else if (column === "transcription") {
            return a.transcription.localeCompare(b.transcription);
        } else if (column === "translation") {
            return a.translation.localeCompare(b.translation);
        } else if (column === "attempt") {
            return a.attempt - b.attempt;
        } else if (column === "guess") {
            return (a.guess === b.guess) ? 0 : a.guess ? -1 : 1;
        }
        return 0;
    });
    updateWordsList();
}

function checkAnswer(card, answer) {
    const cardElement = document.querySelector('.card');

    if (card.translation === answer) {
        currentScoreGuess++;
        card.guess = true;
        cardElement.style.backgroundColor = "green";
        cardElement.style.boxShadow = "0 0 10px green";
    } else {
        card.attempt++;
        cardElement.style.backgroundColor = "red";
        cardElement.style.boxShadow = "0 0 10px red";
    }

    setTimeout(() => {
        cardElement.style.backgroundColor = "";
        cardElement.style.boxShadow = "";
        updateUI();
        
    }, 500);
    
}

function flipCard() {
    const card = document.querySelector('.card');
    card.classList.toggle('flipped');
}

function sendDataToServer() {

    const cardData = {
        data: JSON.parse(JSON.stringify(currentCardsData)),
        scoreGuess: currentScoreGuess,
        scoreAll: currentScoreAll,
        finish: currentFinish
    };

    fetch("/cards", {
            method: "POST",
			headers: {'Content-Type': 'application/json',},
			body: JSON.stringify(cardData),
        })
        .then(response => {
            if (response.ok) {
                document.getElementById("congrats-information").innerText = "The result is fixed!";
                //showAlarm("The result is fixed!");
            }else{
                showAlarm('Возникла проблема с передачей данных.');
            }
            
        })
        .catch(error => {
            showAlarm('Возникла проблема с передачей данных: ' + String(error));
        });
}



