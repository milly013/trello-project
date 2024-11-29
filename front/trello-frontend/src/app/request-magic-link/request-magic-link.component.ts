import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-request-magic-link',
  standalone: true,
  imports: [ReactiveFormsModule, CommonModule, FormsModule,RouterLink],
  templateUrl: './request-magic-link.component.html',
  styleUrl: './request-magic-link.component.css'
})
export class RequestMagicLinkComponent {
  email: string = '';
  message: string = '';
  errorMessage: string = '';

  private apiUrl = 'http://localhost:8000/api/user/users'; // URL backend API-ja

  constructor(private http: HttpClient) {}

  onSubmit() {
    // Send magic link request to backend
    this.http.post(`${this.apiUrl}/request-magic-link`, { email: this.email }).subscribe({
      next: () => {
        this.message = 'Magic link has been successfully sent to your email address.';
        this.errorMessage = '';
      },
      error: (err) => {
        console.error('Error sending magic link:', err);
        this.errorMessage = 'An error occurred while sending the magic link.';
        this.message = '';
      }
    });
  }
}
