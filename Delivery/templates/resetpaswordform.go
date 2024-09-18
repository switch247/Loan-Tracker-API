package templates

// Template for reset password form
const ResetTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password Demo</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }

        .container {
            background-color: white;
            padding: 2rem;
            border-radius: 8px;
            box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
            max-width: 400px;
            width: 100%;
        }

        h1 {
            text-align: center;
            margin-bottom: 2rem;
        }

        .error-message {
            color: red;
            margin-bottom: 1rem;
        }

        label {
            display: block;
            margin-bottom: 0.5rem;
        }

        input {
            width: 100%;
            padding: 0.8rem;
            border: 1px solid #ccc;
            border-radius: 4px;
            font-size: 1rem;
            margin-bottom: 1rem;
        }

        button {
            background-color: #4CAF50;
            color: white;
            border: none;
            padding: 0.8rem 1.5rem;
            border-radius: 4px;
            cursor: pointer;
            font-size: 1rem;
            width: 100%;
        }

        button:hover {
            background-color: #45a049;
        }

        .password-container {
            display: flex;
            align-items: center;
            margin-bottom: 1rem;
            position: relative;
        }

        .password-container input {
            flex-grow: 1;
            margin-bottom: 0;
            padding-right: 2.5rem;
        }

        .show-password-icon {
            position: absolute;
            right: 0.8rem;
            cursor: pointer;
            font-size: 1.2rem;
        }
    </style>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.15.3/css/all.min.css">
</head>
<body>

    <div class="container">
        <h1>Reset Password</h1>
        <form id="reset-form" method="post" action="/users/password-update/{{ .ResetToken }}">
            <div id="error-container" class="error-message"></div>
            <label for="password">New Password:</label>
            <div class="password-container">
                <input type="password" id="password" name="password" required>
                <i class="fas fa-eye show-password-icon"></i>
            </div>
            <button type="submit">Reset Password</button>
        </form>
    </div>

    <script>
        const passwordInput = document.getElementById('password');
        const showPasswordIcon = document.querySelector('.show-password-icon');

        showPasswordIcon.addEventListener('click', () => {
            if (passwordInput.type === 'password') {
                passwordInput.type = 'text';
                showPasswordIcon.classList.remove('fa-eye');
                showPasswordIcon.classList.add('fa-eye-slash');
            } else {
                passwordInput.type = 'password';
                showPasswordIcon.classList.remove('fa-eye-slash');
                showPasswordIcon.classList.add('fa-eye');
            }
        });

        document.getElementById('reset-form').addEventListener('submit', (event) => {
            event.preventDefault();
            const formData = new FormData(event.target);
            fetch(event.target.action, {
                method: 'POST',
                body: formData
            })
            .then(response => {
                if (!response.ok) {
                    return response.json().then(data => {
                        throw new Error(data.message || 'Error resetting password. Please try again.');
                    });
                }
                return response.json();
            })
            .then(data => {
                if (data.message) {
                    alert(data.message);
                    document.getElementById('error-container').textContent = '';
                } else {
                    alert('Password reset successful!');
                }
            })
            .catch(error => {
                document.getElementById('error-container').textContent = error.message;
            });
        });
    </script>
</body>
</html>
`
