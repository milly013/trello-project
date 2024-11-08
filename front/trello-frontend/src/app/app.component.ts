import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { HttpClient } from '@angular/common/http';
import { AppNavComponent } from './app-nav/app-nav.component';

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  imports: [RouterOutlet, CreateProjectComponent, AppNavComponent]
})
export class AppComponent {
  title = 'trello-frontend';
}
