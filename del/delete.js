const loadBtn = document.getElementById('load-btn');
const jsonFilesSelect = document.getElementById('json-files');
const deleteBtn = document.getElementById('delete-btn');
const resultDiv = document.getElementById('result');

loadBtn.addEventListener('click', loadJsonFiles);
deleteBtn.addEventListener('click', deleteJsonFile);

function loadJsonFiles() {
  fetch('26.229.38.10:9000/deletecard')
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
  fetch(`26.229.38.10:9000/deletecard/${selectedFile}`, {
    method: 'DELETE'
  })
  .then(response => response.json())
  .then(_data => {
    resultDiv.innerHTML = `Файл ${selectedFile} удален успешно`;
  })
  .catch(error => console.error('Error:', error));
}