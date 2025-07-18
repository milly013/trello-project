import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../service/auth.service';
import { CommonModule } from '@angular/common';
import { DomSanitizer, SafeValue } from '@angular/platform-browser'; // Import DomSanitizer
import { RECAPTCHA_SETTINGS, RecaptchaModule, RecaptchaSettings } from 'ng-recaptcha';

@Component({
  selector: 'app-login',
  standalone: true,
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
  imports: [ReactiveFormsModule, CommonModule, FormsModule,RouterLink,RecaptchaModule]
})
export class LoginComponent {
  loginForm: FormGroup;
  isSubmitting = false;
  errorMessage: string | null = null;
  siteKey = '6Leoc5EqAAAAAHf_zqSb1gRl6q3YEgigsnVkwZ7I';


  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private router: Router,
    private sanitizer: DomSanitizer
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      recaptcha: [null, Validators.required] 
    });
  }

  onSubmit(): void {
    if (this.loginForm.invalid) {
      console.log("Form is invalid");
      this.errorMessage = 'Invalid input. Please fill out all fields correctly.';
      return;
    }

    const formData = this.loginForm.value;

    // Sanitizacija unosa korisnika kako bi se izbjegao XSS napad
    const sanitizedLoginData = {
      email: this.sanitizeInput(formData.email),
      password: this.sanitizeInput(formData.password),
      recaptchaToken: formData.recaptcha // Send reCAPTCHA token to backend
    };

    console.log("Form is valid, attempting login");
    this.isSubmitting = true;
    this.errorMessage = null;

    this.authService.login(sanitizedLoginData.email, sanitizedLoginData.password, sanitizedLoginData.recaptchaToken).subscribe({
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


  resolved(captchaResponse: string | null): void {
    if (captchaResponse) {
      this.loginForm.get('recaptcha')?.setValue(captchaResponse);
    } else {
      this.loginForm.get('recaptcha')?.setErrors({ required: true });
    }
  }
  
}
