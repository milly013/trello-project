import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink, RouterOutlet } from '@angular/router';
import { AuthService } from '../service/auth.service';
import { HttpClient } from '@angular/common/http';


@Component({
  selector: 'app-app-nav',
  standalone: true,
  imports: [FormsModule, RouterOutlet, RouterLink],
  templateUrl: './app-nav.component.html',
  styleUrls: ['./app-nav.component.css'],

})
export class AppNavComponent {

  constructor(private router: Router){}

  onLogout() {
    localStorage.removeItem('authToken');

    const token = localStorage.getItem('authToken');
    if (!token) {
      this.router.navigate(['/login']);
    }
  }

}
