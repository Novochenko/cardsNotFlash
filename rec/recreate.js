


fetch('https://localhost:443/private/show',{
  method: "GET",
  credentials: "include",
  headers:{
  "Content-Type": "applicatoin/json"
  }
})
.then(response => response.json())
.then(data => {
  const cardsContainer = document.getElementById('cards-container');
  data.forEach(card => {
    const cardElement = document.createElement('div');
    cardElement.innerHTML = `
      <h2>${card.front_side}</h2>
      <p>${card.Back_side}</p>
    `;
    cardsContainer.appendChild(cardElement);
  });
});
// Инициализация
//getJsonFiles();
