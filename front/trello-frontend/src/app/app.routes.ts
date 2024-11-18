import { Routes } from '@angular/router';
import { CreateProjectComponent } from './create-project-component/create-project-component.component';
import { AddTaskComponent } from './add-task/add-task.component';
import { AddUserComponent } from './add-user/add-user.component';
import { RemoveUserComponent } from './remove-user/remove-user.component';
import { RegistrationComponent } from './registration/registration.component';
import { AddMemberToTaskComponent } from './add-member-to-task/add-member-to-task.component';
import { RemoveMemberFromTaskComponent } from './remove-member-from-task/remove-member-from-task.component';
import { VerificationComponent } from './verification/verification.component';
import { HomeComponentComponent } from './home-component/home-component.component';
import { ProjectListComponent } from './project-list/project-list.component';
import { UserListComponent } from './user-list/user-list.component';
import { LoginComponent } from './login/login.component';
import { ProjectDetailsComponent } from './project-details/project-details.component';
import { TaskStatusComponent } from './task-status/task-status.component';


export const appRoutes: Routes = [
  { path: 'project-list/create-project', component: CreateProjectComponent },
  { path: 'add-task',component: AddTaskComponent},
  { path:'add-user',component:AddUserComponent},
  { path:'remove-user'  ,component:RemoveUserComponent},
  { path:'home-page/registration',component:RegistrationComponent},
  { path:'verification',component:VerificationComponent},
  { path:'add-member-to-task',component:AddMemberToTaskComponent},
  { path: 'remove-member-from-task', component: RemoveMemberFromTaskComponent },
  { path: 'home-page', component: HomeComponentComponent},
  { path: 'project-list', component: ProjectListComponent},
  { path: 'user-list', component: UserListComponent},
  { path: 'home-page/login', component: LoginComponent},
  { path: 'project-details/:id', component: ProjectDetailsComponent },
  { path: 'task-status',component:TaskStatusComponent}


];
 