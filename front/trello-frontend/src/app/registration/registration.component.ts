import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http'; // Importuj HttpClient
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-registration',
  templateUrl: './registration.component.html',
  styleUrls: ['./registration.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule] 
})
export class RegistrationComponent {
  registrationForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {
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

    const user = {
      username: formData.username,
      email: formData.email,
      password: formData.password,
      role: formData.isManager ? 'manager' : 'member'
    };

    // Pošaljemo HTTP POST zahtev na /users
    this.http.post('http://localhost:8080/users', user)
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

  get f() {
    return this.registrationForm.controls;
  }
}
