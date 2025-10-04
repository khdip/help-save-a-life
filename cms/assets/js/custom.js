document.getElementById('contactForm').addEventListener('submit', async function (event) {
  event.preventDefault();

  const form = event.target;
  const formData = new FormData(form);

  const nameError = document.getElementById('nameError');
  const emailError = document.getElementById('emailError');
  const commentError = document.getElementById('commentError');
  const formMessage = document.getElementById('formMessage');
  const toastLiveExample = document.getElementById('liveToast');
  const toastBootstrap = bootstrap.Toast.getOrCreateInstance(toastLiveExample);

  // Clear previous messages
  [nameError, emailError, commentError, formMessage].forEach(el => {
    el.innerText = '';
  });

  try {
    const response = await fetch(form.action, {
      method: form.method,
      body: formData,
    });

    const result = await response.json();

    // Handle validation errors
    nameError.innerText = result.FormErrors?.Name || '';
    emailError.innerText = result.FormErrors?.Email || '';
    commentError.innerText = result.FormErrors?.Comment || '';

    if (!response.ok) {
      formMessage.innerText = 'HTTP error! Please try again later.';
      toastLiveExample.classList.remove('bg-success');
      toastLiveExample.classList.add('bg-danger');
      toastBootstrap.show();
      return;
    }

    // Success
    if (result.Message?.SuccessMessage) {
      formMessage.innerText = result.Message.SuccessMessage;
      toastLiveExample.classList.remove('bg-danger');
      toastLiveExample.classList.add('bg-success');
      toastBootstrap.show();
      form.reset();
    }
  } catch (error) {
    formMessage.innerText = 'Server unavailable. The form could not be submitted';
    toastLiveExample.classList.add('bg-danger');
    toastBootstrap.show();
  }
});
