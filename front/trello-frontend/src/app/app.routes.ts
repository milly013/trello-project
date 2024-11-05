import { RouterModule, Routes } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { AppNavComponent } from './app-nav/app-nav.component';
import { AddTaskComponent } from './add-task/add-task.component';

export const appRoutes: Routes = [
  { path: 'create-project', component: CreateProjectComponent },
  { path: 'app-nav', component: AppNavComponent},
  { path: 'add-task', component: AddTaskComponent},
  
  
];
