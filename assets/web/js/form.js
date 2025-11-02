// Employee name autocomplete functionality
document.addEventListener('DOMContentLoaded', function() {
  const input = document.getElementById('employee_name');
  const datalist = document.getElementById('employee-list');
  if (!input || !datalist) return;
  const options = Array.from(datalist.options).map(opt => opt.value);
  let currentFocus = -1;
  input.addEventListener('input', function() {
    const val = this.value.toLowerCase();
    datalist.innerHTML = '';
    options.filter(opt => opt.toLowerCase().includes(val)).forEach(opt => {
      const option = document.createElement('option');
      option.value = opt;
      datalist.appendChild(option);
    });
  });
  input.addEventListener('keydown', function(e) {
    const visibleOptions = Array.from(datalist.options);
    if (e.key === 'ArrowDown') {
      currentFocus = (currentFocus + 1) % visibleOptions.length;
      input.value = visibleOptions[currentFocus]?.value || input.value;
      e.preventDefault();
    } else if (e.key === 'ArrowUp') {
      currentFocus = (currentFocus - 1 + visibleOptions.length) % visibleOptions.length;
      input.value = visibleOptions[currentFocus]?.value || input.value;
      e.preventDefault();
    } else if (e.key === 'Enter') {
      if (currentFocus > -1 && visibleOptions[currentFocus]) {
        input.value = visibleOptions[currentFocus].value;
      }
      currentFocus = -1;
    } else {
      currentFocus = -1;
    }
  });
});
