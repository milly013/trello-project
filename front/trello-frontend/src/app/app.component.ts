import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  imports: [RouterOutlet, CreateProjectComponent,]
})
export class AppComponent {
  title = 'trello-frontend';
}
