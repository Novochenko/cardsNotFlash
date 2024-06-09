fetch('https://localhost:443/private/show',{
  method: "GET",
  credentials: "include",
  headers:{
  "Content-Type": "applicatoin/json"
  }
})
.then(response => response.json())
.then(data => {
  const cardContainer = document.getElementById("edit-card")
  let n = 1;
  data.forEach(card => {
    const cardElement = document.createElement('div');
    cardElement.innerHTML = `
    <div id="edit" class="container-del">
      <div class="card-front">${card.front_side}</div>
      <div class="card-back"><strong>${card.back_side}</strong></div>
    </div>
    `;
    cardElement.dataset.id = n;
    n++;
    cardElement.id = "recreate-id";
    cardContainer.appendChild(cardElement);
  });
  const listItems = cardContainer.children;
  for (let i = 0; i < listItems.length; i++) {
    listItems[i].addEventListener('click', function() {
      const selectedOption = this.dataset.id;

      Swal.fire({
        title: `Введите вопрос и ответ для карточки № ${selectedOption}`,
        html: `
          <label for="question-rec">Вопрос:</label>
          <input id="question-rec" type="text" class="swal2-input">
          <br><br>
          <label for="answer-rec">Ответ:</label>
          <input id="answer-rec" type="text" class="swal2-input">
        `,
        focusConfirm: false,
        showCancelButton: true,
        cancelButtonText: 'Cancel',
        confirmButtonText: 'Submit',
        confirmButtonColor: '#3085d6',
        cancelButtonColor: '#d33',
        allowOutsideClick: false,
        preConfirm: () => {
          const front_side = document.getElementById('question-rec').value;
          const back_side = document.getElementById('answer-rec').value;
          if (!front_side || !back_side) {
            Swal.showValidationMessage('Please enter both question and answer');
          } else {
            // Send data to server or perform other actions
            console.log(`Question: ${front_side}, Answer: ${back_side}`);
            const card_id = parseInt(selectedOption);

            console.log(card_id,"card_id");
            const newCard = {
              "card_id": card_id,
              "front_side": front_side,
              "back_side": back_side
            }
            fetch('https://localhost:443/private/editcard',{
              method: "POST",
              credentials: "include",
              headers:{
                "Content-Type": "aplication/json"
              },
              body:JSON.stringify(newCard)
            })
            .then(response =>{
              if (response.ok){
                console.log('edit success');
                location.reload();
              }
              else{
                console.log('edit fail');
              }
            })
          }
        }
      });

    });
  };
});
