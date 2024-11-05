import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { UserService } from '../service/user.service';
import { ProjectService } from '../service/project.service';
import { CommonModule } from '@angular/common';




@Component({
  selector: 'app-remove-user',
  standalone:true,
  templateUrl: './remove-user.component.html',
  styleUrls: ['./remove-user.component.css'],
  imports: [ReactiveFormsModule,CommonModule,FormsModule],
  providers: [ProjectService]
})


export class RemoveUserComponent implements OnInit {
  removeUserForm: FormGroup;
  users: any[] = []; // Možete definirati tip korisnika prema vašem modelu

  constructor(private fb: FormBuilder, private userService: UserService) {
    this.removeUserForm = this.fb.group({
      userId: ['', Validators.required] // Kontroler za userId
    });
  }

  ngOnInit(): void {
    this.loadUsers(); // Učitajte korisnike prilikom inicijalizacije
  }

  loadUsers() {
    this.userService.getUsers().subscribe((response: any[]) => {
      this.users = response; // Pretpostavljam da response sadrži listu korisnika
    });
  }

  removeUser() {
    if (this.removeUserForm.valid) {
      const userId = this.removeUserForm.get('userId')?.value;
      this.userService.removeUser(userId).subscribe((response: any) => {
        console.log('User removed', response);
        // Ponovo učitajte korisnike nakon brisanja
        this.loadUsers();
        // Opcionalno: Očistite formu ili dodajte poruku o uspehu
        this.removeUserForm.reset();
      });
    }
  }
}
