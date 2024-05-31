let cardList = []


fetch('https://localhost:443/private/show',{
    method: 'GET',
    credentials: 'include',
    headers:{
        'Content-Type': 'application/json',
    },
})
.then(response => {
    if (response.ok){
        console.log("groups on screen");
    }
    else{
        console.error('show error');
    }
return response.json()})
    
.then(data => {
      // Вставляем опции в select
      const cardList = document.getElementById("delete-card")
      data.forEach(card => {
        const cardElement = document.createElement('div');
        cardElement.innerHTML = `
        <div id="box" class="container-del">
          <div class="card-id">${card.id}</div>
          <div class="card-front">${card.front_side}</div>
          <div class="card-back"><strong>${card.back_side}</strong></div>
        </div>
        `;
        cardList.appendChild(cardElement);
        })
})
.catch(error => {
  `Ошибка отправки вопроса: ${error}`;
});
function deleteCard() {
  const deleteButton = document.getElementById('del-btn');
  const selectedOption = document.getElementById('container-list').selectedOptions[0];

  if(selectedOption){
    console.log('rrerere');
    deleteButton.addEventListener('click', (e) => {
      e.preventDefault();
      const id = selectedOption.value;
      const card_id = parseInt(id);
      const cardDel = {
        "card_id": card_id
      }
      fetch('https://localhost:443/private/deletecard',{
        method: 'POST',
        credentials: 'include',
        headers:{
          "Content-Type": "application/json"
        },
        body:JSON.stringify(cardDel)
      })
      .then(response =>{
        if(response.ok){
          console.log('delete success')
        }
        else{
          console.log('delete error')
        }
      })
    })
  }
}
  