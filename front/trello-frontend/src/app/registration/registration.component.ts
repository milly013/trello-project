import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http'; 
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { DomSanitizer, SafeValue } from '@angular/platform-browser';
import { RECAPTCHA_SETTINGS, RecaptchaModule, RecaptchaSettings } from 'ng-recaptcha';


@Component({
  selector: 'app-registration',
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RecaptchaModule]
})
export class RegistrationComponent {
  registrationForm: FormGroup;
  siteKey = '6Leoc5EqAAAAAHf_zqSb1gRl6q3YEgigsnVkwZ7I';


  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router, private sanitizer: DomSanitizer) {
    this.registrationForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      isManager: [false],
      recaptcha: [null, Validators.required]
    });
  }

  register(): void {
    if (this.registrationForm.invalid) {
      alert('Molimo popunite sva obavezna polja ispravno.');
      return;
    }

    const formData = this.registrationForm.value;

    // Sanitizacija unosa korisnika kako bi se izbjegao XSS napad
    const sanitizedUser = {
      username: this.sanitizeInput(formData.username),
      email: this.sanitizeInput(formData.email),
      password: this.sanitizeInput(formData.password),
      role: formData.isManager ? 'manager' : 'member',
      recaptchaToken: formData.recaptcha
    };

    this.http.post('http://localhost:8000/api/user/users', sanitizedUser)
      .subscribe({
        next: (response) => {
          console.log('Korisnik je uspešno registrovan', response);
          alert('Registracija je uspešna!');
          this.router.navigate(['verification']);
        },
        error: (error) => {
          console.error('Greška tokom registracije', error);
          alert('Došlo je do greške prilikom registracije.');
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
    return this.registrationForm.controls;
  }
  resolved(captchaResponse: string | null): void {
    this.registrationForm.patchValue({ recaptcha: captchaResponse });
  }
}
