document.addEventListener('DOMContentLoaded', function () {
    const toastEl = document.querySelector('.toast');
    const toast = bootstrap.Toast.getOrCreateInstance(toastEl);
    toast.show();
});