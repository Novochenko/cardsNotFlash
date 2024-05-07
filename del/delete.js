const loadBtn = document.getElementById('load-btn');
const jsonFilesSelect = document.getElementById('json-files');
const deleteBtn = document.getElementById('delete-btn');
const resultDiv = document.getElementById('result');

loadBtn.addEventListener('click', loadJsonFiles);
deleteBtn.addEventListener('click', deleteJsonFile);

function loadJsonFiles() {
  fetch('localhost:8080/deletecard')
    .then(response => response.json())
    .then(data => {
      const options = data.map(file => {
        return `<option value="${file}">${file}</option>`;
      });
      jsonFilesSelect.innerHTML = options.join('');
    })
    .catch(error => console.error('Error:', error));
}

function deleteJsonFile() {
  const selectedFile = jsonFilesSelect.value;
  fetch(`localhost:8080/deletecard/${selectedFile}`, {
    method: 'DELETE'
  })
  .then(response => response.json())
  .then(data => {
    resultDiv.innerHTML = `Файл ${selectedFile} удален успешно`;
  })
  .catch(error => console.error('Error:', error));
}