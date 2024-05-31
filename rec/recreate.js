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
  data.forEach(card => {
    const cardElement = document.createElement('div');
    cardElement.innerHTML = `
    <div id="edit" class="container-del">
      <div class="card-id">${card.id}</div>
      <div class="card-front">${card.front_side}</div>
      <div class="card-back"><strong>${card.back_side}</strong></div>
    </div>
    `;
    cardContainer.appendChild(cardElement);
  });
});


function editCard(){
  
  const id = document.getElementById('card-id');
  const front_side = document.getElementById('front-side');
  const back_side = document.getElementById('back-side');

  const card_id = parseInt(id);
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

}