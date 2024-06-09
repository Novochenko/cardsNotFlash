let cardList = [];

const addForm = document.getElementById('question-form');

const responseAdd = document.getElementById('response');

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
    document.getElementById('container-list').innerHTML = options.join(' ');
    data.forEach(item => {
        if (!cardList.includes(item.group_id)){
            console.log("progon");
            cardList.push(item.group_id);
            console.log(cardList);
        }
        return data.group_id;
    });
});


function selectChange(){
  const addBtn = document.getElementById('add-btn');
  const selectedOption = document.getElementById('container-list').selectedOptions[0];
  if (selectedOption) {
    addBtn.addEventListener('click', (e) => {
      e.preventDefault();
      const id = selectedOption.value;
      const front_side = document.getElementById('front-side').value;
      const back_side = document.getElementById('back-side').value;
      
      const group_id = parseInt(id);
      console.log(group_id,"group_id");
      console.log(front_side,"   front");
      console.log(back_side,"   back");
      const cardData = {
        "group_id": group_id, 
        "front_side": front_side,
        "back_side": back_side
      }
      // Отправляем обновленный список карточек на сервак
      fetch('https://localhost:443/private/createcard', {
        method: 'POST',
        credentials: "include",
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(cardData)
      })
      .then(response => {
        if (response.ok){
          console.log('good add');
        }
        else{
          console.log('Error: карточка не добавлена');
        }
        return response.json()})
      .then(data => {
        responseDiv.innerHTML = `Вопрос отправлен успешно ID вопроса: ${card_id}`;
        card_id++; // увеличиваем ID вопроса
      })
      .catch(error => {
        `Ошибка отправки вопроса: ${error}`;
      });
    })
    }
  }