// load questionSets into scope
// index 0 will be chosen as default on page load
let questionSetsJSON = [];
let userGroups = [];


fetch('https://localhost:443/private/showallgroups',{
    method: 'GET',
    credentials: 'include',
    headers:{
        'Content-Type': 'application/json',
    },
})
.then(response => {
    if (response.ok){
        console.log("skopiroval group");
    }
    else{
        console.error('ne skopiroval group:error');
    }
return response.json()})
    
.then(data => {
    // Создаем опции для select
    const options = data.map(data => `<option value=${data.group_id}>${data.group_name}</option>`);
    // Вставляем опции в select
    document.getElementById('decks').innerHTML = options.join(' ');
    data.forEach(item => {
        if (!questionSetsJSON.includes(item.group_id)){
            questionSetsJSON.map(item => [item.group_id, item.group_name]);
            console.log("progon");
            userGroups.push(item.group_id);
            console.log(userGroups);
        }
    return console.log(`progon qeustseroatijs ${userGroups[0]}`);
    });
    })

function handleSelectChange() {
    const selectedOption = document.getElementById('decks').selectedOptions[0];
    if (selectedOption) {
      // Получаем значение и текст выбранного элемента
      const group_id = selectedOption.value;
      const name = selectedOption.name;
      const int64 = parseInt (group_id);
      // Выполняем дальнейшие действия с выбранным элементом
      console.log(`Selected option: ${name} (value: ${group_id})`);

      const userData ={
        "group_id": int64
      }
      // Например, можно отправить POST-запрос на сервер с выбранным элементом
      fetch('https://localhost:443/private/showgroupusingtime', {
        method: 'POST',
        credentials: 'include',
        headers:{
        "Content-Type": "application/json"
        },
        body: JSON.stringify(userData)
        })
      .then(response => {
        if (response.ok){
            console.log("gooposd");
        }
        else{
            console.error("error");
        }
       return response.json()})
      .then(data => {
        // data.forEach(questionSetsJSON.map(item => [item.front_side, item.back_side]))
        const cards = data.map(card => ({
            "front_side": card.front_side,
            "back_side": card.back_side
          }));
        const cardshows = document.getElementById('decks');

        cardshows.addEventListener('click', () =>{
            window.addEventListener("beforeunload", (e) => {
            e.preventDefault();
            e.returnValue = "";
        });
        
        // getting the DOM elements 
        // Элементы карточки
        let card = document.querySelector(".card")
        let question = document.querySelector(".question");
        let solution = document.querySelector(".solution");
        let inputField = document.querySelector(".answer");
        // Кнопки
        let checkAnswerBUTTON = document.querySelector(".check_answer");
        let newWordBUTTON = document.querySelector(".new_word");
        let correctBUTTON = document.querySelector(".correct");
        let wrongBUTTON = document.querySelector(".wrong");
        let reloadBUTTON = document.querySelector(".reload");
        let returnBUTTON = document.querySelector(".return");
        // card-deck-choice-fields
        
        // Текст под карточкой
        let remainingCards = document.querySelector(".remaining");
        let knownCards = document.querySelector("#known");
        let knownCardsCounter = 0;
        let nextCards = document.querySelector("#next");

        // define question-set
        // global questionSet
        const questionSet = cards[0];
        let nextRound = [];
        // set a questionSet to start with
        function defineQuestionSet(set) {
            questionSet = set;
        }
        
        // keep track of current deck-index (cardDeckOptions / questionSetsJSON)
        let lastDeckIDX = 0;
        // load the deck the user selects from the options-drop-down
        let deckOptions = document.querySelector("#decks");
        deckOptions.addEventListener("change", (e) => {
            let selectedDeck = deckOptions.value - 1;
            let lastDeckIDX = cardDeckOptions.indexOf(selectedDeck);
        
            defineQuestionSet(questionSetsJSON[lastDeckIDX]);
            newCard();
            });
        
        // повернуть сторону back_sideа на front_side
        returnBUTTON.addEventListener("click", () => card.classList.remove("flipped"));
        
        
        // взять рандомную пару (front_side/back_side) из всех
        function getQuestionPair(dict) {
            let rand = Math.floor(Math.random() * dict.length);
            return Object.entries(dict)[rand][1];   
        }
        
        // innitial randomPair on Page-Load
        let randomPair = getQuestionPair(cards);
        
        // display first question
        displayQuestion(randomPair);
        
        // write question on the front of card
        function displayQuestion(rP) {
            // get curser into input-field
            if (inputField) inputField.focus(); 
            // remove input field from page if it is not needed (e.g. for rechtquest)
            if (rP["input"]) {
                inputField.classList.remove("hidden");
                wrongBUTTON.classList.add("hidden");
                correctBUTTON.classList.add("hidden");
            } else {
                inputField.classList.add("hidden");
                wrongBUTTON.classList.remove("hidden");
                correctBUTTON.classList.remove("hidden");
            }
            // turn card to front-side
            card.classList.remove("flipped");
            // write question on front-side of the card
            question.innerHTML = rP["front_side"];
            // add hidden to the last answer, so the card-size rescales down (to question-size)
            solution.classList.add("hidden");
        
            // display current stack of cards
            remainingCards.innerHTML = `Есть еще ${cards.length} карт в колоде`
            knownCards.innerHTML = `Известно на данный момент: ${knownCardsCounter}`; 
            nextCards.innerHTML = `Следующий раунд: ${nextRound.length}`; 
        }
        
        // used inside of "flipBackAndDisplayAnswer" to split multiple answers
        function splitPhraseIfSeveralNumbers(phrase) {
            let re = /\d\.\s.+\;/;
            if (re.test(phrase)) {
                phrase = phrase.split(";");
            }
            return phrase;
        }
        
        // flip card to back-side and display/render answer(s)
        function flipBackAndDisplayAnswer() {
            // display answer on the back of the card
            card.classList.add("flipped");
            // the solution-text is hidden (so the card-size isn't too big from the last answer)
            solution.classList.remove("hidden");
        
            // if no input no "nächste Frage"
            if (randomPair["input"]) newWordBUTTON.classList.remove("hidden");
            else newWordBUTTON.classList.add("hidden");
        
            // create backside of the card
            let answer = document.querySelector(".answer");
            if (randomPair["input"]) {
                answer = answer.value;
                if (answer == randomPair["back_side"]) solution.innerHTML = "Правильно!";
                else solution.innerHTML = `К сожалению нет.<br>Правильный ответ был <em>"${randomPair["back_side"]}"</em>.`;
            } else {
                // create List of possible multiple-answer
                let answerList = splitPhraseIfSeveralNumbers(randomPair["back_side"]);
                // if it is just one answer display it
                if (typeof answerList == "string") solution.innerHTML = randomPair["back_side"];
                // if muslitple answers display them as a list
                else {
                    solution.innerHTML = "";
                    let answerListDOM = document.createElement("ol");
                    solution.appendChild(answerListDOM);
                    answerList.forEach(a => {
                        // check if a number is in front and delete it so the ordered list tag provides numbers;
                        let numRegEx = /^\s*\d+\.\s/g;
                        a = a.replace(numRegEx, "");
                        let listElement = document.createElement("li");
                        answerListDOM.appendChild(listElement);
                        listElement.innerHTML = a;
                    });
                }
            } 
        
            // clear the input field if there is one for next questions
            let answerInput = document.querySelector(".answer");
            if (answerInput) answerInput.value = "";
        }
        
        // get a new card
        function newCard() {
            if (questionSet.length === 0 && nextRound.length > 0) {
                questionSet = nextRound;
                questionSetsJSON[lastDeckIDX] = nextRound;
                nextRound = [];
            }
            // create new randomPair in global scope
            randomPair = getQuestionPair(cards);
            displayQuestion(randomPair);
        }
        
        // removes current randomPair of question Answer from global questionSet-Array of objects
        function removeCardFromSet(correct) {
            let idx = cards.findIndex(qa => qa["front_side"] == randomPair["front_side"]);
            let card = cards[idx];
            if (correct) {
                cards.splice(idx, 1);
                knownCardsCounter += 1;
            } else {
                nextRound.push(card);
                cards.splice(idx, 1);
            }
            if (cards.length > 0 || nextRound.length > 0) newCard();
        }
        
        // attache the show-result function to the button on frontside of card
        checkAnswerBUTTON.addEventListener("click", flipBackAndDisplayAnswer);
        
        // attache newWord-function to button on backside of card
        newWordBUTTON.addEventListener("click", newCard);
        
        // attache functionality "newCard" to wrong-button
        wrongBUTTON.addEventListener("click", () => {
            removeCardFromSet(false);
        });
        correctBUTTON.addEventListener("click", () => {
            removeCardFromSet(true);
        });
        
        // reload the page / begin from the beginning
        reloadBUTTON.addEventListener("click", () => location.reload());}
        )
      }
    )
    }
}