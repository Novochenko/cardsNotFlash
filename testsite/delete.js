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
      let k = 1;
      const cardList = document.getElementById("delete-card"); 
      data.forEach(card => {
        const cardElement = document.createElement('div');
        cardElement.innerHTML = `
        <div id="box" class="container-del">
          <div class="card-front">${card.front_side}</div>
          <div class="card-back"><strong>${card.back_side}</strong></div>
        </div>
        `;
        cardElement.className = "delete-id";
        cardElement.dataset.id = card.id;
        cardElement.id = k;
        k++;
        //cardElement.classList = "delete-list"
        cardList.appendChild(cardElement);
      });
      const listItems = cardList.children;
      for (let i = 0; i < listItems.length; i++) {
        listItems[i].addEventListener('click', function() {
          // Выбрали элемент, теперь можно взаимодействовать с ним
          const selectedOption = this.dataset.id;
          
          Swal.fire({
            title: `Выберите элемент с ID ${selectedOption}`,
            text: 'Вы уверены, что хотите выбрать этот элемент?',
            icon: 'question',
            showCancelButton: true,
            confirmButtonText: 'Да, выбрать',
            cancelButtonText: 'Отмена'
          })
          .then((result =>{
            if (result.isConfirmed){
              console.log(`Выбран элемент №: ${selectedOption}`);
        
              if(selectedOption){
                console.log('rrerere');
                const id = selectedOption;
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
                    location.reload();
                    console.log('delete success')
                    
                  }
                  else{
                    console.log('delete error')
                  }
                })
              }
            }
            else{
              console.log('Действие отменено');
            }
          }))


})
}
})