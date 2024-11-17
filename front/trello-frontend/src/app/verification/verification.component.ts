import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-verification',
  templateUrl: './verification.component.html',
  styleUrls: ['./verification.component.css'],
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule]
})
export class VerificationComponent {
  verificationForm: FormGroup;

  constructor(private fb: FormBuilder, private http: HttpClient, private router: Router) {
    this.verificationForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      verificationCode: ['', [Validators.required]]
    });
  }

  verify(): void {
    if (this.verificationForm.invalid) {
      alert('Please enter all required fields correctly.');
      return;
    }

    const formData = this.verificationForm.value;

    // Dinamičko kreiranje URL-a sa vrednostima iz forme
    const url = `http://localhost:8080/verify/${formData.email}/${formData.verificationCode}`;

    // HTTP POST request for verification
    this.http.post(url, {})
      .subscribe({
        next: (response) => {
          console.log('Verification successful', response);
          alert('You have successfully verified your account!');
          this.router.navigate(['home-page'])
        },
        error: (error) => {
          console.error('Verification error', error);
          alert('An error occurred during verification.');
        }
      });
  }
}
