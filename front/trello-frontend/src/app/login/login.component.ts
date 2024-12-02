import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../service/auth.service';
import { CommonModule } from '@angular/common';
import { DomSanitizer, SafeValue } from '@angular/platform-browser'; // Import DomSanitizer

declare const grecaptcha: any;

@Component({
  selector: 'app-login',
  standalone: true,
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  imports: [ReactiveFormsModule, CommonModule, FormsModule, RouterLink]
})
export class LoginComponent implements OnInit {
  loginForm: FormGroup;
  isSubmitting = false;
  errorMessage: string | null = null;
  grecaptchaReady = false;

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private router: Router,
    private sanitizer: DomSanitizer
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      captchaToken: ['']
    });
  }

  ngOnInit(): void {
    // Dinamičko učitavanje reCAPTCHA skripte
    const script = document.createElement('script');
    script.src = 'https://www.google.com/recaptcha/enterprise.js?render=6Lesmo8qAAAAAEcDN2Qb5c6uKOnZSARfvn8DA_8v';
    script.async = true;
    script.defer = true;
    script.onload = () => {
      if (typeof grecaptcha !== 'undefined') {
        grecaptcha.ready(() => {
          console.log('reCAPTCHA is ready');
          this.grecaptchaReady = true;
        });
      } else {
        console.error('grecaptcha is not defined');
      }
    };
    document.head.appendChild(script);
  }

  onSubmit(): void {
    if (this.loginForm.invalid) {
      console.log("Form is invalid");
      this.errorMessage = 'Invalid input. Please fill out all fields correctly.';
      return;
    }

    if (!this.grecaptchaReady) {
      console.error('grecaptcha is not ready');
      this.errorMessage = 'Captcha is not ready yet. Please try again later.';
      return;
    }

    const formData = this.loginForm.value;

    // Sanitizacija unosa korisnika kako bi se izbjegao XSS napad
    const sanitizedLoginData = {
      email: this.sanitizeInput(formData.email),
      password: this.sanitizeInput(formData.password),
      captchaToken: ''
    };

    console.log("Form is valid, attempting login");
    this.isSubmitting = true;
    this.errorMessage = null;

    grecaptcha.ready(() => {
      console.log('reCAPTCHA is ready for execution');
      grecaptcha.execute('6Lesmo8qAAAAAEcDN2Qb5c6uKOnZSARfvn8DA_8v', { action: 'login' })
        .then((token: string) => {
          sanitizedLoginData.captchaToken = token;

          console.log("Form is valid, attempting login with reCAPTCHA");

          this.authService.login(sanitizedLoginData.email, sanitizedLoginData.password, sanitizedLoginData.captchaToken).subscribe({
            next: (response) => {
              // Čuvanje JWT tokena u local storage
              console.log("Login successful, saving token");
              localStorage.setItem('authToken', response.token);
              localStorage.setItem('userRole', response.userRole);
              this.router.navigate(['home-page']);
            },
            error: (error) => {
              console.error("Login failed", error);
              this.errorMessage = 'Invalid email or password';
              this.isSubmitting = false;
            }
          });
        }).catch((error: any) => {
          console.error("reCAPTCHA execution failed", error);
          this.errorMessage = 'Captcha verification failed. Please try again.';
          this.isSubmitting = false;
        });
    });
  }

  sanitizeInput(input: string): string {
    // Escape-ovanje opasnih znakova kako bi se spriječili XSS napadi
    const element: HTMLElement = document.createElement('div');
    element.innerText = input;
    return element.innerHTML;
  }

  get f() {
    return this.loginForm.controls;
  }
}
