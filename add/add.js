

const addForm = document.getElementById('question-form');
const addBtn = document.getElementById('add-btn');
const responseAdd = document.getElementById('response');

addBtn.addEventListener('click', (e) => {
  e.preventDefault();
  const group_id = document.getElementById('group_id').value;
  const front_side = document.getElementById('question').value;
  const back_side = document.getElementById('answer').value;

  fetch('https://localhost:443/createcard')
    .then(response => response.json())
    .then(data => {
      const selectElement = document.getElementById('groups');
      data.forEach(group => {
        const optionElement = document.createElement('option');
        optionElement.text = group.name; // или любое другое поле из объекта group
        optionElement.value = group.group_id;
        selectElement.appendChild(optionElement);
      });
    })
    .catch(error => console.error('Error:', error));



  // Создаем объект вопроса и ответа
fetch('https://localhost:443/createcard')
.then(response => response.json())
.then(createcard => {
  const card_id = createcard.length > 0 ? createcard[createcard.length - 1].id + 1 : 1;
  const cardData = {
    group: group_id, 
    id: card_id,
    question: front_side,
    answer: back_side
  };
  createcard.push(cardData);

  // Отправляем обновленный список карточек на сервак
  fetch('https://localhost:443/createcard', {
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