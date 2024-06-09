// // Add event listener to sidebar links
// document.querySelectorAll('.sidebar a').forEach(link => {
//     link.addEventListener('click', event => {
//         event.preventDefault();
//         console.log(`Link clicked: ${link.textContent}`);
//     });
// });

// Получаем элементы рабочих пространств
// const userWorkspace = document.getElementById('user-workspace');
// const homeWorkspace = document.getElementById('home-workspace');
// const editCardWorkspace = document.getElementById('edit-card-workspace');
// const addCardWorkspace = document.getElementById('add-card-workspace');
// const deleteCardWorkspace = document.getElementById('delete-card-workspace');
// const authorsWorkspace = document.getElementById('authors-workspace');

// Получить все ссылки на названия рабочих областей
const links = document.querySelectorAll('a');

// Функция для отображения рабочей области
function showArea(areaId) {
  // Скрыть все рабочие области
  document.querySelectorAll('.area').forEach(area => area.classList.remove('.active'));
  
  // Отобразить выбранную рабочую область
  document.getElementById(areaId).classList.add('.active');
}

// Добавить обработчик события для каждой ссылки
links.forEach(link => {
  link.addEventListener('click', event => {
    event.preventDefault();
    const areaId = link.getAttribute('data-area');
    showArea(areaId);
  });
});

// document.addEventListener("DOMContentLoaded", function(event) {
//   var element = document.querySelector(".start-screen");
//   element.classList.add("active");
// });

// // Функция для перехода между рабочими пространствами
// function switchWorkspace(workspaceId) {
//   // Скрыть все рабочие пространства
//   userWorkspace.style.display = 'none';
//   homeWorkspace.style.display = 'none';
//   editCardWorkspace.style.display = 'none';
//   addCardWorkspace.style.display = 'none';
//   deleteCardWorkspace.style.display = 'none';
//   authorsWorkspace.style.display = 'none';

//   // Показать выбранное рабочее пространство
//   const workspace = document.getElementById(workspaceId);
//   workspace.classList.add('visible');
//   setTimeout(() => {
//     workspace.classList.remove('hidden');
//   }, 500);
//   //workspace.style.display = 'block';
// }

// // Добавляем обработчики событий для кнопок перехода


// document.getElementById('home-button').addEventListener('click', () => {
//   switchWorkspace('home-workspace');
// });

// document.getElementById('user-button').addEventListener('click', () => {
//   switchWorkspace('user-workspace');
// });

// document.getElementById('add-card-button').addEventListener('click', () => {
//   switchWorkspace('add-card-workspace');
// });

// document.getElementById('edit-card-button').addEventListener('click', () => {
//   switchWorkspace('edit-card-workspace');
// });

// document.getElementById('delete-card-button').addEventListener('click', () => {
//   switchWorkspace('delete-card-workspace');
// });

// document.getElementById('authors-button').addEventListener('click', () => {
//   switchWorkspace('authors-workspace');
// });