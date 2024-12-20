import { Component } from '@angular/core';
import { RouterModule, RouterOutlet } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { HttpClient } from '@angular/common/http';
import { AppNavComponent } from './app-nav/app-nav.component';

import { RemoveUserComponent } from './remove-user/remove-user.component';
import { AddTaskComponent } from './add-task/add-task.component';
import { AddUserComponent } from './add-user/add-user.component';
import { HomeComponentComponent } from './home-component/home-component.component';
import { FormsModule } from '@angular/forms';


@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  imports: [RouterOutlet,AppNavComponent,FormsModule]

})
export class AppComponent {
  title = 'trello-frontend';
}
