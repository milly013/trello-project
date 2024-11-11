import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http'; // Importuj HttpClient
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-registration',
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule]  // Direktno uvozimo potrebne module
})
export class RegistrationComponent {
  registrationForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient) { // Dodaj HttpClient u konstruktor
    this.registrationForm = this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      confirmPassword: ['', [Validators.required]]
    });
  }

  register(): void {
    if (this.registrationForm.invalid) {
      alert('Please complete all required fields correctly.');
      return;
    }

    const formData = this.registrationForm.value;

    // Pošaljemo HTTP POST zahtev na /users
    this.http.post('http://localhost:8080/users', formData)
      .subscribe({
        next: (response) => {
          console.log('User registered successfully', response);
          alert('Registration successful!');
        },
        error: (error) => {
          console.error('Error during registration', error);
          alert('An error occurred while registering.');
        }
      });
  }

  get f() {
    return this.registrationForm.controls;
  }
}
