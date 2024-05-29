let cardList = [];

const addForm = document.getElementById('question-form');

const responseAdd = document.getElementById('response');

// addBtn.addEventListener('click', (e) => {
//   e.preventDefault();
//   const group_id = document.getElementById('group_id').value;
//   const front_side = document.getElementById('question').value;
//   const back_side = document.getElementById('answer').value;

//   // const userData = {
//   //   front_side,
//   //   back_side,
//   //   group_id
//   // };

//   // fetch('https://localhost:443/private/createcard', {
//   //   method: "POST",
//   //   headers: {
//   //     "Content-Type": "application/json"
//   //   },
//   //   body: JSON.stringify(userData)
//   // })
//   //   .then((response) =>
//   //     {      
//   //       if (response.ok){
//   //         console.log('card was created')
//   //       }
//   //       else{
//   //         console.log('error to add card')
//   //       }
//   //       return response.json()
//   //     } 

//   //   )
//   //   .then(data => {
//   //     const selectElement = document.getElementById('groups');
//   //     data.forEach(group => {
//   //       const optionElement = document.createElement('option');
//   //       optionElement.text = group.name; // или любое другое поле из объекта group
//   //       optionElement.value = group.group_id;
//   //       selectElement.appendChild(optionElement);
//   //     });
//   //   })
//   //   .catch(error => console.error('Error:', error));



//   // Создаем объект вопроса и ответа
// fetch('https://localhost:443/private/createcard')
// .then(response => response.json())
// .then(createcard => {
//   const card_id = createcard.length > 0 ? createcard[createcard.length - 1].id + 1 : 1;
//   const cardData = {
//     "group_id": group_id, 
//     "front_side": front_side,
//     "back_side": back_side
//   };
//   createcard.push(cardData);

//   // Отправляем обновленный список карточек на сервак
//   fetch('https://localhost:443/private/createcard', {
//     method: 'POST',
//     headers: {
//       'Content-Type': 'application/json'
//     },
//     body: JSON.stringify(cardData)
//   })
//   .then(response => {
//     if (response.ok){
//       console.log('good add');
//     }
//     else{
//       console.log('Error: карточка не добавлена');
//     }
//     return response.json()})
//   .then(_data => {
//     responseDiv.innerHTML = `Вопрос отправлен успешно ID вопроса: ${card_id}`;
//     card_id++; // увеличиваем ID вопроса
//   })
//   .catch(error => {
//     responseDiv.innerHTML = `Ошибка отправки вопроса: ${error}`;
//   });
// })
// });

const cardLists = document.getElementById('fetch-cards');

//const currentGroup = document.getElementById('choose-group');

cardLists.addEventListener('click', () => {
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
})

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
      const cardData = {
        "group_id": group_id, 
        "front_side": front_side,
        "back_side": back_side
      }
      // Отправляем обновленный список карточек на сервак
      fetch('https://localhost:443/private/createcard', {
        method: 'POST',
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