import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http'; 
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { DomSanitizer, SafeValue } from '@angular/platform-browser';

@Component({
  selector: 'app-registration',
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule]
})
export class RegistrationComponent {
  registrationForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router, private sanitizer: DomSanitizer) {
    this.registrationForm = this.fb.group({
      username: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      isManager: [false]
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
      role: formData.isManager ? 'manager' : 'member'
    };

    this.http.post('http://localhost:8080/users', sanitizedUser)
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
}
