import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AuthService } from '../service/auth.service';

@Component({
  selector: 'app-home-component',
  standalone: true,
  imports: [RouterLink, CommonModule],
  templateUrl: './home-component.component.html',
  styleUrl: './home-component.component.css'
})
export class HomeComponentComponent {
  constructor(public authService: AuthService) { }

  isLoggedIn(){
    return this.authService.isLoggedIn();
  }
}
