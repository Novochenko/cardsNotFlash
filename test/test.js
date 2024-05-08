const form = document.getElementById('registration-form');

form.addEventListener('submit', (e) => {
  e.preventDefault();

  // const username = document.getElementById('username').value;
  // const email = document.getElementById('email').value;
  // const password = document.getElementById('password').value;
  const test = document.getElementById('password').value;
  // const userData = {
  //   username,
  //   email,
  //   password
  // };
  const userData = {
    test
  };

  fetch('http://127.0.0.1:9000/test', {
    method: 'POST',
    mode: 'no-cors',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(userData)
  })
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error(error));
});