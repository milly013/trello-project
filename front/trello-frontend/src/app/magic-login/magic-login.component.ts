import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-magic-login',
  standalone: true,
  imports: [ReactiveFormsModule, CommonModule, FormsModule,RouterLink],
  templateUrl: './magic-login.component.html',
  styleUrl: './magic-login.component.css'
})
export class MagicLoginComponent implements OnInit {
  token: string = '';
  errorMessage: string = '';

  private apiUrl = 'http://localhost:8000/api/user/users'; // URL backend API-ja

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private http: HttpClient
  ) {}

  ngOnInit(): void {
    // Get the token from URL if available for magic login
    this.route.queryParams.subscribe(params => {
      this.token = params['token'];
      if (this.token) {
        this.magicLogin();
      }
    });
  }

  magicLogin() {
    console.log('Attempting to log in with token:', this.token);
    // Send token to backend for verification
    this.http.post(`${this.apiUrl}/magic-login`, { token: this.token }).subscribe({
      next: (response: any) => {
        console.log('Login successful, response:', response);
        // Store JWT token and user role in localStorage or session
        localStorage.setItem('authToken', response.token);
        localStorage.setItem('userRole', response.userRole);
        // Redirect user to home page
        this.router.navigate(['/home-page']).then(success => {
          if (!success) {
            this.errorMessage = 'Failed to navigate to home page.';
          }
        });
      },
      error: (err) => {
        console.error('Login error:', err);
        this.errorMessage = 'An error occurred during the login process.';
      }
    });
  }
}
