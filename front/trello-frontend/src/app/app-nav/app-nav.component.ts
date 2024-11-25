
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink, RouterOutlet } from '@angular/router';
import { AuthService } from '../service/auth.service';
import { HttpClient } from '@angular/common/http';
import { CommonModule } from '@angular/common';
import { UserService } from '../service/user.service';



@Component({
  selector: 'app-app-nav',
  standalone: true,
  imports: [FormsModule, RouterOutlet, RouterLink, CommonModule],
  templateUrl: './app-nav.component.html',
  styleUrls: ['./app-nav.component.css'],

})
export class AppNavComponent {

  constructor(private router: Router,){}

  onLogout() {
    localStorage.removeItem('authToken');
    localStorage.removeItem('managerId');
    localStorage.removeItem('userRole');

    const token = localStorage.getItem('authToken');
    if (!token) {
      this.router.navigate(['home-page/login']);
    }
  }
  isLoggedIn(): boolean {
    const token = localStorage.getItem('authToken');
    if (token) {
      return true;
    }
    return false
  }


}
