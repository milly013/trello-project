import { Component } from '@angular/core';
import { RouterModule, RouterOutlet } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { HttpClient } from '@angular/common/http';
import { AppNavComponent } from './app-nav/app-nav.component';
import { AddTaskComponent } from './add-task/add-task.component';
import { appRoutes } from './app.routes';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  imports: [
    RouterOutlet, 
    AppNavComponent,
    FormsModule
  ]
})
export class AppComponent {
  title = 'trello-frontend';
}
