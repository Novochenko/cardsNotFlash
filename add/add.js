const addForm = document.getElementById('question-form');
const addBtn = document.getElementById('add-btn');
const responseAdd = document.getElementById('response');

addBtn.addEventListener('click', (e) => {
  e.preventDefault();
  const front_side = document.getElementById('question').value;
  const back_side = document.getElementById('answer').value;

  // Создаем объект вопроса и ответа
fetch('127.0.0.1:9000/createcard')
.then(response => response.json())
.then(createcard => {
  const card_id = createcard.length > 0 ? createcard[createcard.length - 1].id + 1 : 1;
  const cardData = {
    id: card_id,
    question: front_side,
    answer: back_side
  };
  createcard.push(cardData);

  // Отправляем обновленный список карточек на сервак
  fetch('127.0.0.1:9000/createcard', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(cardData)
  })
  .then(response => {
    if (response.ok){
      console.log('good');
    }
    else{
      console.log('Error: карточка не добавлена');
    }
    return response.json()})
  .then(_data => {
    responseDiv.innerHTML = `Вопрос отправлен успешно ID вопроса: ${card_id}`;
    card_id++; // увеличиваем ID вопроса
  })
  .catch(error => {
    responseDiv.innerHTML = `Ошибка отправки вопроса: ${error}`;
  });
})
});