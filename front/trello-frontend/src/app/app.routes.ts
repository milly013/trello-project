import { Routes } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { AddTaskComponent } from './add-task/add-task.component';

export const appRoutes: Routes = [
  { path: '', component: CreateProjectComponent },
  { path: 'add-task', component: AddTaskComponent}
  
];
