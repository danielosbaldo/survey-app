// Reset form after successful employee creation
document.addEventListener('htmx:afterSwap', function(evt) {
    if (evt.detail.target.id === 'employees-content') {
        const form = document.getElementById('employee-form');
        if (form) {
            form.reset();
            // Reset the shops dropdown to initial state
            const container = document.getElementById('employee-shops-dropdown-container');
            if (container) {
                container.innerHTML = `
                    <select name="shop_ids" multiple class="w-full p-3 border border-gray-300 rounded-xl text-sm focus:ring-2 focus:ring-pink-400">
                        <option value="">Primero selecciona una ciudad</option>
                    </select>
                `;
            }
        }
    }
});

// Debug when this script loads
console.log('employees_section script loaded');

// Debug when section is rendered
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded - checking for employee container');
    const container = document.getElementById('employee-shops-dropdown-container');
    console.log('Container found on DOM load?', container !== null);
});

// Simple HTMX debugging
document.addEventListener('htmx:targetError', function(evt) {
    console.error('HTMX Target Error - target:', evt.detail.target);
    console.error('Element that triggered:', evt.detail.elt);

    // Check if our specific target exists globally
    const container = document.getElementById('employee-shops-dropdown-container');
    console.log('employee-shops-dropdown-container exists globally?', container !== null);

    // Check if target exists within the form
    const form = evt.detail.elt.closest('form');
    const containerInForm = form ? form.querySelector('#employee-shops-dropdown-container') : null;
    console.log('employee-shops-dropdown-container exists in form?', containerInForm !== null);

    // Log all elements with similar IDs for debugging
    const allContainers = document.querySelectorAll('[id*="dropdown-container"]');
    console.log('All dropdown containers found:', allContainers.length);
    allContainers.forEach((el, i) => console.log(`Container ${i}:`, el.id, el));

    // Log the entire form HTML for debugging
    console.log('Form HTML:', form ? form.innerHTML : 'No form found');
});

document.addEventListener('htmx:beforeRequest', function(evt) {
    console.log('HTMX request starting to:', evt.detail.requestConfig?.path);
    console.log('Request target:', evt.detail.target);

    // Verify target exists before request
    const targetId = evt.detail.target;
    if (targetId && typeof targetId === 'string' && targetId.startsWith('#')) {
        const targetEl = document.querySelector(targetId);
        console.log('Target element found before request?', targetEl !== null);
    }
});

document.addEventListener('htmx:afterRequest', function(evt) {
    console.log('HTMX request completed, status:', evt.detail.xhr?.status);
});

// Debug when employees section is loaded via HTMX
document.addEventListener('htmx:afterSettle', function(evt) {
    console.log('HTMX afterSettle event');
    const container = document.getElementById('employee-shops-dropdown-container');
    console.log('Container exists after settle?', container !== null);

    // Re-process HTMX for dynamic content
    if (evt.detail.elt.querySelector && evt.detail.elt.querySelector('#employee-shops-dropdown-container')) {
        console.log('Processing HTMX for dynamically loaded employee section');
        htmx.process(evt.detail.elt);
    }
});
