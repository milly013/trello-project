import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { MagicLinkService } from '../service/magic-link.service';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-magic-link',
  standalone: true,
  imports: [CommonModule,ReactiveFormsModule,FormsModule],
  templateUrl: './magic-link.component.html',
  styleUrl: './magic-link.component.css'
})
export class MagicLinkComponent implements OnInit {
  email: string = '';
  token: string = '';
  message: string = '';
  errorMessage: string = '';
  http: any;

  constructor(
    private magicLinkService: MagicLinkService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    // Get the token from URL if available for magic login
    this.route.queryParams.subscribe(params => {
      this.token = params['token'];
      if (this.token && this.token.trim() !== '') {
        this.magicLogin();
      } else {
        this.errorMessage = 'Token not found or is invalid.';
      }
    });
  }

  onSubmit() {
    // Koristi MagicLinkService za slanje zahteva
    this.magicLinkService.requestMagicLink(this.email).subscribe({
      next: () => {
        this.message = 'Magic link has been successfully sent to your email address.';
        this.errorMessage = '';
      },
      error: () => {
        this.errorMessage = 'An error occurred while sending the magic link.';
        this.message = '';
      }
    });
  }

  magicLogin() {
    // Send token to backend for verification
    this.http.post('/api/magic-login', { token: this.token }).subscribe({
      next: (response: any) => {
        // Store JWT token and user role in localStorage or session
        localStorage.setItem('authToken', response.token);
        localStorage.setItem('userRole', response.userRole);
        // Redirect user to home page
        this.router.navigate(['/home-page']);
      },
      error: () => {
        this.errorMessage = 'An error occurred during the login process.';
      }
    });
  }
}